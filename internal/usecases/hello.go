package usecases

import (
	"errors"
	"github.com/hashicorp/go-multierror"
	. "github.com/krajeswaran/gostartup/internal/adapters"
	"github.com/krajeswaran/gostartup/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"gopkg.in/eapache/go-resiliency.v1/breaker"
	"strconv"
)

// HelloRepo contains HelloRepo usecases
type HelloRepo struct{}

//FetchUserName fetches user name from DB, given user Id
func (h *HelloRepo) FetchUserName(userId string) (string, error) {
	if userId == "" {
		// log here, because these errors are not propagated upwards
		log.Err(h.updateApiStats(true)).Msg("updating API stats")
		return "", errors.New("passed empty user Id")
	}

	user, err := DB.FetchUser(userId)
	if err != nil {
		log.Err(h.updateApiStats(true)).Msg("updating API stats")
		return "", multierror.Append(err, errors.New("failed to fetch user for id: "+userId))
	}

	log.Err(h.updateApiStats(false)).Msg("updating API stats")
	return user.Name, nil
}

func (h *HelloRepo) updateApiStats(didApiFail bool) error {
	// check feature flag first
	if !viper.GetBool("stats.ff_enable_stats") {
		log.Debug().Msg("feature flag turned OFF for stats")
		return nil
	}

	// update stats in cache with a circuit breaker
	b := breaker.New(viper.GetInt("stats.error_threshold"), viper.GetInt("stats.success_threshold"), viper.GetDuration("stats.cache_timeout"))
	b.Go(func() error {
		count, err := Cache.UpdateApiStats(didApiFail)
		if err != nil {
			return multierror.Append(err, errors.New("failed to update API stats in cache"))
		}

		// after a certain api hits reset stats counter
		if count >= viper.GetInt64("stats.reset_stats_on_api_count") {
			log.Err(Cache.ResetApiStats()).Msg("resetting API stats")
		}

		return nil
	})

	return nil
}

//GetApiStats Fetches API stats from cache
func (h *HelloRepo) GetApiStats() (*models.Stat, error) {
	statsRaw, err := Cache.GetApiStats()
	if err != nil {
		return nil, multierror.Append(err, errors.New("failed to fetch API stats from cache"))
	}

	stat := models.Stat{}
	var newErr error
	stat.ApiTotalCount, err = strconv.ParseUint(statsRaw[0], 10, 64)
	stat.ApiFailureCount, newErr = strconv.ParseUint(statsRaw[1], 10, 64)

	err = multierror.Append(err, newErr)
	if err != nil {
		return nil, err
	}

	return &stat, nil
}

//CreateUser Creates a user in DB
func (h *HelloRepo) CreateUser(name string) (*models.User, error) {
	user, err := DB.CreateUser(name)
	if err != nil {
		return nil, multierror.Append(err, errors.New("failed to create user in DB"))
	}

	return user, nil
}
