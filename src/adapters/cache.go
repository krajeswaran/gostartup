package adapters

import (
	"errors"
	"gostartup/src/config"
	"gostartup/src/models"
	"strconv"
	"time"
	"github.com/shopspring/decimal"
	"github.com/satori/go.uuid"
)

// TODO extract field consts
// CacheAdapter - Struct to logically bind all the cache related functions
type CacheAdapter struct{}

const NL_CACHE_PREFIX = "nl:"
const ROUTELIST_CACHE_PREFIX = "rl:"
const ROUTE_CACHE_PREFIX = "r:"
const PROVIDER_CACHE_PREFIX = "p:"
const STAT_CACHE_PREFIX = "s:"
const MDR_CACHE_PREFIX = "m:"
const BALANCE_CACHE_PREFIX = "b:"
const STAT_CACHE_SL_PREFIX = "s:sw:"
const STAT_CACHE_API_SUCCESS_RATIO = "api_suc_ratio"
const STAT_CACHE_API_LATENCY_RATIO = "api_lat_ratio"
const STAT_CACHE_DELIVERY_SUCCESS_RATIO = "del_suc_ratio"
const STAT_CACHE_DELIVERY_LATENCY_RATIO = "del_lat_ratio"
const STAT_CACHE_DLR_CODES = "dlr_codes"

// DeepStatus checks health of redis
func (c *CacheAdapter) DeepStatus() error {
	if err := Cache().Ping().Err(); err != nil {
		return errors.New("SERVICE_CACHE_DOWN")
	}
	return nil
}

func (c *CacheAdapter) SetMccMncForNumber(number string, mcc string, mnc string) error {
	values := map[string]interface{}{
		"mcc": mcc,
		"mnc": mnc,
	}
	_, err := Cache().HMSet(NL_CACHE_PREFIX+number, values).Result()

	return err
}

func (c *CacheAdapter) GetMccMncForNumber(number string) (string, string, error) {
	vals, err := Cache().HMGet(NL_CACHE_PREFIX+number, "mcc", "mnc").Result()
	if err != nil {
		return "", "", err
	}

	mcc, ok := vals[0].(string)
	if !ok {
		return "", "", errors.New("Nil value for mcc")
	}

	mnc, ok := vals[1].(string)
	if !ok {
		return "", "", errors.New("Nil value for mnc")
	}

	return mcc, mnc, nil
}

func (c *CacheAdapter) GetProvider(providerId uint) (*models.Provider, error) {
	providerRaw, err := Cache().Get(PROVIDER_CACHE_PREFIX + strconv.Itoa(int(providerId))).Result()
	if err != nil {
		return nil, err
	}
	provider := models.Provider{}
	_, err = provider.UnmarshalMsg([]byte(providerRaw))
	if err != nil {
		return nil, err
	}

	return &provider, nil
}

func (c *CacheAdapter) DeleteProvider(providerId uint) error {
	_, err := Cache().Del(PROVIDER_CACHE_PREFIX + strconv.Itoa(int(providerId))).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheAdapter) GetRoutesByMccMnc(mcc string, mnc string) ([]models.SmsRoute, error) {
	vals, err := Cache().LRange(ROUTELIST_CACHE_PREFIX+mcc+mnc, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(vals) == 0 {
		return nil, errors.New("routelist not found for mcc/mnc")
	}

	var routes []models.SmsRoute
	// in case of any cache miss, abort the whole thing and let db fill up
	for _, routeId := range vals {
		// fetch each route id
		routeRaw, err := Cache().Get(ROUTE_CACHE_PREFIX + routeId).Result()
		if err != nil {
			return nil, err
		}
		route := models.SmsRoute{}
		cost := models.SmsCost{}
		_, err = cost.UnmarshalMsg([]byte(routeRaw))
		if err != nil {
			return nil, err
		}
		route.SmsCost = cost

		// fetch provider for each route id
		provider, err := c.GetProvider(route.ProviderId)
		if err != nil {
			return nil, err
		}
		route.Provider = *provider

		// fetch stats
		stat, err := c.GetSmsStatForRoute(routeId)
		if err != nil {
			return nil, err
		}
		route.SmsStat = *stat

		// append super struct back
		routes = append(routes, route)
	}

	return routes, nil
}

func (c *CacheAdapter) SetProvider(providerId uint, provider *models.Provider) error {
	providerRaw, err := provider.MarshalMsg(nil)
	if err != nil {
		return err
	}
	_, err = Cache().Set(PROVIDER_CACHE_PREFIX+strconv.Itoa(int(providerId)), providerRaw,
		DefaultCacheExpiry).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheAdapter) SetRoutesByMccMnc(mcc string, mnc string, routes []models.SmsRoute) error {
	for _, route := range routes {
		if err := c.SetProvider(route.ProviderId, &route.Provider); err != nil {
			return err
		}

		routeRaw, err := route.SmsCost.MarshalMsg(nil)
		if err != nil {
			return err
		}
		routeIdStr := strconv.Itoa(int(route.SmsCost.ID))
		_, err = Cache().Set(ROUTE_CACHE_PREFIX+routeIdStr, routeRaw,
			DefaultCacheExpiry).Result()
		if err != nil {
			return err
		}

		// set stats
		values := map[string]interface{}{
			"success_rate": route.SmsStat.SuccessRateComputed,
			"latency": route.SmsStat.LatencyComputed,
			"api_success_count": route.SmsStat.ApiSuccessCount,
			"api_total_count": route.SmsStat.ApiTotalCount,
			"delivery_success_count": route.SmsStat.DeliverySuccessCount,
			"delivery_total_count": route.SmsStat.DeliveryTotalCount,
			"manual_override": route.SmsStat.ManualOverride,
		}
		_, err = Cache().HMSet(STAT_CACHE_PREFIX+routeIdStr, values).Result()
		if err != nil {
			return err
		}

		// add route.id to mcc/mnc list
		_, err = Cache().LPush(ROUTELIST_CACHE_PREFIX+mcc+mnc, route.SmsCost.ID).Result()
		if err != nil {
			return err
		}
	}

	// all done
	return nil
}

func (c *CacheAdapter) GetSmsStatForRoute(routeId string) (*models.SmsStat, error) {
	stat := models.SmsStat{}
	var err error
	raw, err := Cache().HMGet(STAT_CACHE_PREFIX+routeId, "success_rate", "latency", "api_success_count",
		"api_total_count", "delivery_success_count", "delivery_total_count", "manual_override").Result()
	if err != nil {
		return nil, err
	}
	smsCostIdInt, _ := strconv.ParseUint(routeId, 10, 32)
	stat.SmsCostId = uint(smsCostIdInt)
	stat.SuccessRateComputed, _ = strconv.ParseFloat(raw[0].(string), 64)
	stat.LatencyComputed, _ = strconv.ParseFloat(raw[1].(string), 64)
	stat.ApiSuccessCount, _ = strconv.ParseUint(raw[2].(string), 10, 64)
	stat.ApiTotalCount, _ = strconv.ParseUint(raw[3].(string), 10, 64)
	stat.DeliverySuccessCount, _ = strconv.ParseUint(raw[4].(string), 10, 64)
	stat.DeliveryTotalCount, _ = strconv.ParseUint(raw[5].(string), 10, 64)
	stat.ManualOverride, _ = strconv.ParseBool(raw[6].(string))

	// all done
	return &stat, nil
}

func (c *CacheAdapter) UpdateApiStats(routeId string, isApiSuccess bool, latency float64) error {
	totalCount, err := Cache().HIncrBy(STAT_CACHE_PREFIX+routeId, "api_total_count", 1).Result()
	if err != nil {
		return err
	}

	if isApiSuccess {
		// TODO multi-exec this shit
		successCount, err := Cache().HIncrBy(STAT_CACHE_PREFIX+routeId, "api_success_count", 1).Result()
		if err != nil {
			return err
		}

		var successRatio float64
		successRatio = float64(successCount / totalCount)

		if err := slidingWindowUpdate(STAT_CACHE_SL_PREFIX + routeId + STAT_CACHE_API_SUCCESS_RATIO,
			successRatio); err != nil {
			return err
		}
		if err := slidingWindowUpdate(STAT_CACHE_SL_PREFIX + routeId + STAT_CACHE_API_LATENCY_RATIO,
			latency); err != nil {
			return err
		}
	}

	return nil
}

func (c *CacheAdapter) UpdateDeliveryStats(routeId string, isDeliverySuccess bool, latency float64,
	dlrError string) error {
	totalCount, err := Cache().HIncrBy(STAT_CACHE_PREFIX+routeId, "delivery_total_count", 1).Result()
	if err != nil {
		return err
	}

	if isDeliverySuccess {
		// TODO multi-exec this shit
		successCount, err := Cache().HIncrBy(STAT_CACHE_PREFIX+routeId, "delivery_success_count", 1).Result()
		if err != nil {
			return err
		}

		var successRatio float64
		successRatio = float64(successCount / totalCount)

		if err := slidingWindowUpdate(STAT_CACHE_SL_PREFIX + routeId + STAT_CACHE_DELIVERY_SUCCESS_RATIO,
			successRatio); err != nil {
			return err
		}
		if err := slidingWindowUpdate(STAT_CACHE_SL_PREFIX + routeId + STAT_CACHE_DELIVERY_LATENCY_RATIO,
			latency); err != nil {
			return err
		}
	}

	if err := slidingWindowUpdate(STAT_CACHE_SL_PREFIX + routeId + STAT_CACHE_DLR_CODES,
		dlrError); err != nil {
		return err
	}
	// compute computed stats only *after* delivery
	apiWeighted, err := calculateListAvg(STAT_CACHE_SL_PREFIX+routeId+STAT_CACHE_API_SUCCESS_RATIO,
		config.GetConfig().GetFloat64("stat.api_success_ratio_weightage"))
	if err != nil {
		return err
	}

	apiLatencyWeighted, err := calculateListAvg(STAT_CACHE_SL_PREFIX+routeId+STAT_CACHE_API_LATENCY_RATIO,
		config.GetConfig().GetFloat64("stat.api_latency_ratio_weightage"))
	if err != nil {
		return err
	}

	// repeat for delivery stuff
	delWeighted, err := calculateListAvg(STAT_CACHE_SL_PREFIX+routeId+STAT_CACHE_DELIVERY_SUCCESS_RATIO,
		config.GetConfig().GetFloat64("stat.delivery_success_ratio_weightage"))
	if err != nil {
		return err
	}

	delLatencyWeighted, err := calculateListAvg(STAT_CACHE_SL_PREFIX+routeId+STAT_CACHE_DELIVERY_LATENCY_RATIO,
		config.GetConfig().GetFloat64("stat.delivery_latency_ratio_weightage"))
	if err != nil {
		return err
	}

	// now compute the computed
	// (0.88 * 0.80-del weight + 0.98 * 0.20-api weight) = success compute
	successComputed := delWeighted + apiWeighted
	latencyComputed := delLatencyWeighted + apiLatencyWeighted

	values := map[string]interface{}{
		"success_rate": successComputed,
		"latency": latencyComputed,
	}
	_, err = Cache().HMSet(STAT_CACHE_PREFIX+routeId, values).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheAdapter) UpdateMDR(mdr *models.MessageDetailRecord) error {
	mdrRaw, err := mdr.MarshalMsg(nil)
	if err != nil {
		return err
	}
	mdrExpiry := time.Duration(config.GetConfig().GetInt("sync.mdr_interval")) * time.Minute
	_, err = Cache().Set(MDR_CACHE_PREFIX+mdr.ID.String(), mdrRaw, mdrExpiry).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheAdapter) GetMDR(msgId string) (*models.MessageDetailRecord, error) {
	mdrRaw, err := Cache().Get(MDR_CACHE_PREFIX + msgId).Result()
	if err != nil {
		return nil, err
	}
	mdr := models.MessageDetailRecord{}
	_, err = mdr.UnmarshalMsg([]byte(mdrRaw))
	if err != nil {
		return nil, err
	}

	return &mdr, nil
}

func slidingWindowUpdate(key string, val interface{}) error {
	windowLen, err := Cache().RPush(key, val).Result()
	if err != nil {
		return err
	}

	maxWindowLen := config.GetConfig().GetInt64("stat.sliding_window_max_length")
	if windowLen > maxWindowLen {
		_, err := Cache().LTrim(key, 0, windowLen - maxWindowLen).Result()
		if err != nil {
			return err
		}
	}

	return nil
}

func calculateListAvg(key string, weight float64) (float64, error) {
	result, err := Cache().LRange(key, 0, -1).Result()
	if err != nil {
		return -1, nil
	}
	var sum float64
	for _, r := range result {
		ratio, _ := strconv.ParseFloat(r, 64)
		sum += ratio
	}
	avg := sum / float64(len(result))
	weightedAvg := weight * avg
	return weightedAvg, nil
}

func (c *CacheAdapter) GetBalance(acctId uuid.UUID) (*models.Balance, error) {
	return c.GetBalanceFromKey(BALANCE_CACHE_PREFIX + acctId.String())
}

func (c *CacheAdapter) GetBalanceFromKey(key string) (*models.Balance, error) {
	balanceRaw, err := Cache().HGetAll(key).Result()
	if err != nil {
		return nil, err
	}
	balance := models.Balance{}
	balance.Balance, err = decimal.NewFromString(balanceRaw["bal"])
	if err != nil {
		return nil, err
	}

	balance.ThresholdAmount, err = decimal.NewFromString(balanceRaw["threshold"])
	if err != nil {
		return nil, err
	}

	// since only enabled items are cached
	balance.IsEnabled = true

	if err != nil {
		return nil, err
	}

	return &balance, nil
}

func (c *CacheAdapter) UpdateBalance(acctId uuid.UUID, amount decimal.Decimal) (decimal.Decimal, decimal.Decimal, error) {
	acctIdStr := acctId.String()
	ok, _ := Cache().Exists(BALANCE_CACHE_PREFIX + acctIdStr).Result()
	if ok < 1 {
		return decimal.Zero, decimal.Zero, errors.New("balance not found")
	}

	pipe := Cache().TxPipeline()
	amountFloat, _ := amount.Round(6).Float64()
	updated := pipe.HIncrByFloat(BALANCE_CACHE_PREFIX + acctIdStr, "bal", amountFloat)
	threshold := pipe.HGet(BALANCE_CACHE_PREFIX + acctIdStr, "threshold")

	_, err := pipe.Exec()
	if err != nil {
		return decimal.Zero, decimal.Zero, err
	}
	updatedAmt := decimal.NewFromFloat(updated.Val())
	thresholdAmt, _ := decimal.NewFromString(threshold.Val())

	return updatedAmt, thresholdAmt, nil
}

func (c *CacheAdapter) WriteBalance(acctId uuid.UUID, balance *models.Balance) error {
	values := map[string]interface{}{
		"bal": balance.Balance.Round(6).String(),
		"threshold": balance.ThresholdAmount.Round(6).String(),
	}
	_, err := Cache().HMSet(BALANCE_CACHE_PREFIX + acctId.String(), values).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheAdapter) DeleteBalance(acctId uuid.UUID, balance *models.Balance) error {
	_, err := Cache().Del(BALANCE_CACHE_PREFIX+acctId.String()).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *CacheAdapter) SetLock(key string) error {
	lockExpiry := time.Duration(config.GetConfig().GetInt("sync.lock_expiry")) * time.Minute
	ok, err := Cache().SetNX(key, "1", lockExpiry).Result()
	if ok {
		return nil
	} else {
		return errors.New(key + " lock already exists, error: " + err.Error())
	}
}

func (c *CacheAdapter) ReleaseLock(key string) error {
	_, err := Cache().Del(key).Result()
	if err != nil {
		return err
	}

	return err
}

func (c *CacheAdapter) PullBalances() ([]models.Balance, error) {
	var cursor uint64
	var keys []string
	for {
		var err error
		keys, cursor, err = Cache().Scan(cursor, BALANCE_CACHE_PREFIX, 100).Result()
		if err != nil {
			return nil, err
		}
		if cursor == 0 {
			break
		}
	}

	// now pull each value
	var balances []models.Balance
	for _, k := range keys {
		bal, err := c.GetBalanceFromKey(k)
		if err != nil {
			return nil, err
		}
		balances = append(balances, *bal)
	}

	return balances, nil
}


func (c *CacheAdapter) PullStats() ([]models.SmsStat, error) {
	var cursor uint64
	var keys []string
	for {
		var err error
		keys, cursor, err = Cache().Scan(cursor, STAT_CACHE_PREFIX, 100).Result()
		if err != nil {
			return nil, err
		}
		if cursor == 0 {
			break
		}
	}

	// now pull each value
	var stats []models.SmsStat
	for _, k := range keys {
		stat, err := c.GetSmsStatForRoute(k)
		if err != nil {
			return nil, err
		}
		stats = append(stats, *stat)
	}

	return stats, nil
}
