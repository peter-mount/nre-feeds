package darwinkb

import "time"

// All of the feeds require a security token from this endPoint
//
// authenticateMaxAge the time we allow a token to exist before refreshing it.
//                    Tokens last for an hour but we keep them for a shorter period to ensure that
//                    they remain valid - just in case the clocks on both ends are out of sync.
// authenticateUrl    the endpoint url
//
// see setToken in retriever.go
const (
	authenticateMaxAge = 45 * time.Minute
	authenticateUrl    = "https://opendata.nationalrail.co.uk/authenticate"
)

// Details about the remote xml rest endpoints:
//
// Bucket   the bucket used for this feed
// MaxAge   the max age for data before performing an update, used to reduce unnecessary calls
// Schedule the Cron schedule for this feed
// Url      the remote URL for this feed
const (
	// Incidents updated every 15 minutes except during the early hours
	incidentsBucket   = "incidents"
	incidentsMaxAge   = 9 * time.Minute
	incidentsSchedule = "0 0/15 0-1,5-23 * * *"
	incidentsUrl      = "https://opendata.nationalrail.co.uk/api/staticfeeds/5.0/incidents"

	// Public Promotions - currently not implemented
	publicPromotionsUrl = "https://opendata.nationalrail.co.uk/api/staticfeeds/4.0/promotions-publics"

	// National Service Indicators every 10 minutes except during the early hours
	serviceIndicatorsBucket   = "serviceIndicators"
	serviceIndicatorsMaxAge   = 9 * time.Minute
	serviceIndicatorsSchedule = "0 0/10 0-1,4-23 * * *"
	serviceIndicatorsUrl      = "https://opendata.nationalrail.co.uk/api/staticfeeds/4.0/serviceIndicators"

	// Stations updated once an hour except during the early hours
	stationsBucket   = "stations"
	stationsMaxAge   = 2 * time.Hour
	stationsSchedule = "0 30 4-9 * * *"
	stationsUrl      = "https://opendata.nationalrail.co.uk/api/staticfeeds/4.0/stations"

	// Ticket Restrictions - currently not implemented
	ticketRestrictionsUrl = "https://opendata.nationalrail.co.uk/api/staticfeeds/4.0/ticket-restrictions"

	// Ticket types once an hour during the morning only as only updated infrequently
	ticketTypesBucket   = "ticketTypes"
	ticketTypesMaxAge   = 6 * time.Hour
	ticketTypesSchedule = "0 40 4-9 * * *"
	ticketTypesUrl      = "https://opendata.nationalrail.co.uk/api/staticfeeds/4.0/ticket-types"

	// Tocs/Companies once an hour during the morning only as only updated infrequently
	tocsBucket   = "companies"
	tocsMaxAge   = time.Hour
	tocsSchedule = "0 35 4-9 * * *"
	tocsUrl      = "https://opendata.nationalrail.co.uk/api/staticfeeds/4.0/tocs"
)
