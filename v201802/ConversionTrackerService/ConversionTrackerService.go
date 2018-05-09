package ConversionTrackerService

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

//
// Type of snippet code to generate.
//
type AdWordsConversionTrackerTrackingCodeType string

const (

	//
	// The snippet that is fired as a result of a website page loading.
	//
	AdWordsConversionTrackerTrackingCodeTypeWEBPAGE AdWordsConversionTrackerTrackingCodeType = "WEBPAGE"

	//
	// The snippet contains a JavaScript function which fires the tag. This function is typically
	// called from an onClick handler added to a link or button element on the page.
	//
	AdWordsConversionTrackerTrackingCodeTypeWEBPAGE_ONCLICK AdWordsConversionTrackerTrackingCodeType = "WEBPAGE_ONCLICK"

	//
	// For embedding on a (mobile) webpage. The snippet contains a JavaScript function which fires
	// the tag. This function is typically called from an onClick handler added to a link or button
	// element on the page that also instructs a mobile device to dial the advertiser's phone
	// number.
	//
	AdWordsConversionTrackerTrackingCodeTypeCLICK_TO_CALL AdWordsConversionTrackerTrackingCodeType = "CLICK_TO_CALL"
)

type AppConversionAppConversionType string

const (
	AppConversionAppConversionTypeNONE AppConversionAppConversionType = "NONE"

	AppConversionAppConversionTypeDOWNLOAD AppConversionAppConversionType = "DOWNLOAD"

	AppConversionAppConversionTypeIN_APP_PURCHASE AppConversionAppConversionType = "IN_APP_PURCHASE"

	AppConversionAppConversionTypeFIRST_OPEN AppConversionAppConversionType = "FIRST_OPEN"
)

//
// App platform for the AppConversionTracker.
//
type AppConversionAppPlatform string

const (
	AppConversionAppPlatformNONE AppConversionAppPlatform = "NONE"

	AppConversionAppPlatformITUNES AppConversionAppPlatform = "ITUNES"

	AppConversionAppPlatformANDROID_MARKET AppConversionAppPlatform = "ANDROID_MARKET"

	AppConversionAppPlatformMOBILE_APP_CHANNEL AppConversionAppPlatform = "MOBILE_APP_CHANNEL"
)

type AppPostbackUrlErrorReason string

const (

	//
	// Invalid Url format.
	//
	AppPostbackUrlErrorReasonINVALID_URL_FORMAT AppPostbackUrlErrorReason = "INVALID_URL_FORMAT"

	//
	// Invalid domain.
	//
	AppPostbackUrlErrorReasonINVALID_DOMAIN AppPostbackUrlErrorReason = "INVALID_DOMAIN"

	//
	// Some of the required macros were not found.
	//
	AppPostbackUrlErrorReasonREQUIRED_MACRO_NOT_FOUND AppPostbackUrlErrorReason = "REQUIRED_MACRO_NOT_FOUND"
)

//
// Attribution models describing how to distribute credit for a particular
// conversion across potentially many prior interactions. See
// https://support.google.com/adwords/answer/6259715 for more information about
// attribution modeling in AdWords.
//
type AttributionModelType string

const (
	AttributionModelTypeUNKNOWN AttributionModelType = "UNKNOWN"

	//
	// Attributes all credit for a conversion to its last click.
	//
	AttributionModelTypeLAST_CLICK AttributionModelType = "LAST_CLICK"

	//
	// Attributes all credit for a conversion to its first click.
	//
	AttributionModelTypeFIRST_CLICK AttributionModelType = "FIRST_CLICK"

	//
	// Attributes credit for a conversion equally across all of its clicks.
	//
	AttributionModelTypeLINEAR AttributionModelType = "LINEAR"

	//
	// Attributes exponentially more credit for a conversion to its more recent clicks
	// (half-life is 1 week).
	//
	AttributionModelTypeTIME_DECAY AttributionModelType = "TIME_DECAY"

	//
	// Attributes 40% of the credit for a conversion to its first and last clicks.
	// Remaining 20% is evenly distributed across all other clicks.
	//
	AttributionModelTypeU_SHAPED AttributionModelType = "U_SHAPED"

	//
	// Flexible model that uses machine learning to determine the appropriate
	// distribution of credit among clicks.
	//
	AttributionModelTypeDATA_DRIVEN AttributionModelType = "DATA_DRIVEN"
)

//
// The single reason for the authentication failure.
//
type AuthenticationErrorReason string

const (

	//
	// Authentication of the request failed.
	//
	AuthenticationErrorReasonAUTHENTICATION_FAILED AuthenticationErrorReason = "AUTHENTICATION_FAILED"

	//
	// Client Customer Id is required if CustomerIdMode is set to CLIENT_EXTERNAL_CUSTOMER_ID.
	// Starting version V201409 ClientCustomerId will be required for all requests except
	// for {@link CustomerService#get}
	//
	AuthenticationErrorReasonCLIENT_CUSTOMER_ID_IS_REQUIRED AuthenticationErrorReason = "CLIENT_CUSTOMER_ID_IS_REQUIRED"

	//
	// Client Email is required if CustomerIdMode is set to CLIENT_EXTERNAL_EMAIL_FIELD.
	//
	AuthenticationErrorReasonCLIENT_EMAIL_REQUIRED AuthenticationErrorReason = "CLIENT_EMAIL_REQUIRED"

	//
	// Client customer Id is not a number.
	//
	AuthenticationErrorReasonCLIENT_CUSTOMER_ID_INVALID AuthenticationErrorReason = "CLIENT_CUSTOMER_ID_INVALID"

	//
	// Client customer Id is not a number.
	//
	AuthenticationErrorReasonCLIENT_EMAIL_INVALID AuthenticationErrorReason = "CLIENT_EMAIL_INVALID"

	//
	// Client email is not a valid customer email.
	//
	AuthenticationErrorReasonCLIENT_EMAIL_FAILED_TO_AUTHENTICATE AuthenticationErrorReason = "CLIENT_EMAIL_FAILED_TO_AUTHENTICATE"

	//
	// No customer found for the customer id provided in the header.
	//
	AuthenticationErrorReasonCUSTOMER_NOT_FOUND AuthenticationErrorReason = "CUSTOMER_NOT_FOUND"

	//
	// Client's Google Account is deleted.
	//
	AuthenticationErrorReasonGOOGLE_ACCOUNT_DELETED AuthenticationErrorReason = "GOOGLE_ACCOUNT_DELETED"

	//
	// Google account login token in the cookie is invalid.
	//
	AuthenticationErrorReasonGOOGLE_ACCOUNT_COOKIE_INVALID AuthenticationErrorReason = "GOOGLE_ACCOUNT_COOKIE_INVALID"

	//
	// problem occurred during Google account authentication.
	//
	AuthenticationErrorReasonFAILED_TO_AUTHENTICATE_GOOGLE_ACCOUNT AuthenticationErrorReason = "FAILED_TO_AUTHENTICATE_GOOGLE_ACCOUNT"

	//
	// The user in the google account login token does not match the UserId in the cookie.
	//
	AuthenticationErrorReasonGOOGLE_ACCOUNT_USER_AND_ADS_USER_MISMATCH AuthenticationErrorReason = "GOOGLE_ACCOUNT_USER_AND_ADS_USER_MISMATCH"

	//
	// Login cookie is required for authentication.
	//
	AuthenticationErrorReasonLOGIN_COOKIE_REQUIRED AuthenticationErrorReason = "LOGIN_COOKIE_REQUIRED"

	//
	// User in the cookie is not a valid Ads user.
	//
	AuthenticationErrorReasonNOT_ADS_USER AuthenticationErrorReason = "NOT_ADS_USER"

	//
	// Oauth token in the header is not valid.
	//
	AuthenticationErrorReasonOAUTH_TOKEN_INVALID AuthenticationErrorReason = "OAUTH_TOKEN_INVALID"

	//
	// Oauth token in the header has expired.
	//
	AuthenticationErrorReasonOAUTH_TOKEN_EXPIRED AuthenticationErrorReason = "OAUTH_TOKEN_EXPIRED"

	//
	// Oauth token in the header has been disabled.
	//
	AuthenticationErrorReasonOAUTH_TOKEN_DISABLED AuthenticationErrorReason = "OAUTH_TOKEN_DISABLED"

	//
	// Oauth token in the header has been revoked.
	//
	AuthenticationErrorReasonOAUTH_TOKEN_REVOKED AuthenticationErrorReason = "OAUTH_TOKEN_REVOKED"

	//
	// Oauth token HTTP header is malformed.
	//
	AuthenticationErrorReasonOAUTH_TOKEN_HEADER_INVALID AuthenticationErrorReason = "OAUTH_TOKEN_HEADER_INVALID"

	//
	// Login cookie is not valid.
	//
	AuthenticationErrorReasonLOGIN_COOKIE_INVALID AuthenticationErrorReason = "LOGIN_COOKIE_INVALID"

	//
	// Failed to decrypt the login cookie.
	//
	AuthenticationErrorReasonFAILED_TO_RETRIEVE_LOGIN_COOKIE AuthenticationErrorReason = "FAILED_TO_RETRIEVE_LOGIN_COOKIE"

	//
	// User Id in the header is not a valid id.
	//
	AuthenticationErrorReasonUSER_ID_INVALID AuthenticationErrorReason = "USER_ID_INVALID"
)

//
// The reasons for the authorization error.
//
type AuthorizationErrorReason string

const (

	//
	// Could not complete authorization due to an internal problem.
	//
	AuthorizationErrorReasonUNABLE_TO_AUTHORIZE AuthorizationErrorReason = "UNABLE_TO_AUTHORIZE"

	//
	// Customer has no AdWords account.
	//
	AuthorizationErrorReasonNO_ADWORDS_ACCOUNT_FOR_CUSTOMER AuthorizationErrorReason = "NO_ADWORDS_ACCOUNT_FOR_CUSTOMER"

	//
	// User doesn't have permission to access customer.
	//
	AuthorizationErrorReasonUSER_PERMISSION_DENIED AuthorizationErrorReason = "USER_PERMISSION_DENIED"

	//
	// Effective user doesn't have permission to access this customer.
	//
	AuthorizationErrorReasonEFFECTIVE_USER_PERMISSION_DENIED AuthorizationErrorReason = "EFFECTIVE_USER_PERMISSION_DENIED"

	//
	// Access denied since the customer is not active.
	//
	AuthorizationErrorReasonCUSTOMER_NOT_ACTIVE AuthorizationErrorReason = "CUSTOMER_NOT_ACTIVE"

	//
	// User has read-only permission cannot mutate.
	//
	AuthorizationErrorReasonUSER_HAS_READONLY_PERMISSION AuthorizationErrorReason = "USER_HAS_READONLY_PERMISSION"

	//
	// No customer found.
	//
	AuthorizationErrorReasonNO_CUSTOMER_FOUND AuthorizationErrorReason = "NO_CUSTOMER_FOUND"

	//
	// Developer doesn't have permission to access service.
	//
	AuthorizationErrorReasonSERVICE_ACCESS_DENIED AuthorizationErrorReason = "SERVICE_ACCESS_DENIED"
)

//
// Enums for the various reasons an error can be thrown as a result of
// ClientTerms violation.
//
type ClientTermsErrorReason string

const (

	//
	// Customer has not agreed to the latest AdWords Terms & Conditions
	//
	ClientTermsErrorReasonINCOMPLETE_SIGNUP_CURRENT_ADWORDS_TNC_NOT_AGREED ClientTermsErrorReason = "INCOMPLETE_SIGNUP_CURRENT_ADWORDS_TNC_NOT_AGREED"
)

//
// Conversion deduplication mode for Conversion Optimizer. That is, whether to
// optimize for number of clicks that get at least one conversion, or total number
// of conversions per click.
//
type ConversionDeduplicationMode string

const (

	//
	// Number of clicks that get at least one conversion.
	//
	ConversionDeduplicationModeONE_PER_CLICK ConversionDeduplicationMode = "ONE_PER_CLICK"

	//
	// Total number of conversions per click.
	//
	ConversionDeduplicationModeMANY_PER_CLICK ConversionDeduplicationMode = "MANY_PER_CLICK"
)

//
// The category of conversion tracker that is being tracked.
//
type ConversionTrackerCategory string

const (
	ConversionTrackerCategoryDEFAULT ConversionTrackerCategory = "DEFAULT"

	ConversionTrackerCategoryPAGE_VIEW ConversionTrackerCategory = "PAGE_VIEW"

	ConversionTrackerCategoryPURCHASE ConversionTrackerCategory = "PURCHASE"

	ConversionTrackerCategorySIGNUP ConversionTrackerCategory = "SIGNUP"

	ConversionTrackerCategoryLEAD ConversionTrackerCategory = "LEAD"

	ConversionTrackerCategoryREMARKETING ConversionTrackerCategory = "REMARKETING"

	//
	// Download is applicable only to {@link AppConversion} types,
	// and is an error to use in conjunction with other types.
	// AppConversions must use download only if they also specify
	// {@link AppConversion#appConversionType} as DOWNLOAD or FIRST_OPEN.
	// If any other appConversionType is used, then some other category besides
	// DOWNLOAD must be used.
	//
	ConversionTrackerCategoryDOWNLOAD ConversionTrackerCategory = "DOWNLOAD"
)

//
// Status of the conversion tracker. The user cannot ADD or SET the
// status to {@code HIDDEN}.
//
type ConversionTrackerStatus string

const (

	//
	// Visits to the conversion page will be recorded.
	//
	ConversionTrackerStatusENABLED ConversionTrackerStatus = "ENABLED"

	//
	// Visits to the conversion page will not be recorded.
	//
	ConversionTrackerStatusDISABLED ConversionTrackerStatus = "DISABLED"

	//
	// Conversions will be recorded, but the conversion tracker will not appear in the UI.
	//
	ConversionTrackerStatusHIDDEN ConversionTrackerStatus = "HIDDEN"
)

//
// Enumerates all the possible reasons for a ConversionTypeError.
//
type ConversionTrackingErrorReason string

const (

	//
	// An attempt to make a forked conversion type from a global conversion type was made,
	// but there already exists a conversion type forked from this global conversion type.
	//
	ConversionTrackingErrorReasonALREADY_CREATED_CUSTOM_CONVERSION_TYPE ConversionTrackingErrorReason = "ALREADY_CREATED_CUSTOM_CONVERSION_TYPE"

	//
	// This user is not whitelisted for the import of Analytics goals and profiles, and yet
	// requested to mutate an Analytics conversion type.
	//
	ConversionTrackingErrorReasonANALYTICS_NOT_ALLOWED ConversionTrackingErrorReason = "ANALYTICS_NOT_ALLOWED"

	//
	// Cannot execute an ADD operation on this subclass of ConversionType (currently, only
	// instances of AdWordsConversionType may be added).
	//
	ConversionTrackingErrorReasonCANNOT_ADD_CONVERSION_TYPE_SUBCLASS ConversionTrackingErrorReason = "CANNOT_ADD_CONVERSION_TYPE_SUBCLASS"

	//
	// Creating an upload conversion type with isExternallyAttributedConversion and
	// isSalesforceConversion both set is not allowed.
	//
	ConversionTrackingErrorReasonCANNOT_ADD_EXTERNALLY_ATTRIBUTED_SALESFORCE_CONVERSION ConversionTrackingErrorReason = "CANNOT_ADD_EXTERNALLY_ATTRIBUTED_SALESFORCE_CONVERSION"

	//
	// AppConversions cannot change app conversion types once it has been set.
	//
	ConversionTrackingErrorReasonCANNOT_CHANGE_APP_CONVERSION_TYPE ConversionTrackingErrorReason = "CANNOT_CHANGE_APP_CONVERSION_TYPE"

	//
	// AppConversions cannot change app platforms once it has been set.
	//
	ConversionTrackingErrorReasonCANNOT_CHANGE_APP_PLATFORM ConversionTrackingErrorReason = "CANNOT_CHANGE_APP_PLATFORM"

	//
	// Cannot change between subclasses of ConversionType
	//
	ConversionTrackingErrorReasonCANNNOT_CHANGE_CONVERSION_SUBCLASS ConversionTrackingErrorReason = "CANNNOT_CHANGE_CONVERSION_SUBCLASS"

	//
	// If a conversion type's status is initially non-hidden, it may not be changed to Hidden;
	// nor may hidden conversion types be created through the API. Hidden conversion types are
	// typically created by backend processes.
	//
	ConversionTrackingErrorReasonCANNOT_SET_HIDDEN_STATUS ConversionTrackingErrorReason = "CANNOT_SET_HIDDEN_STATUS"

	//
	// The user attempted to change the Category when it was uneditable.
	//
	ConversionTrackingErrorReasonCATEGORY_IS_UNEDITABLE ConversionTrackingErrorReason = "CATEGORY_IS_UNEDITABLE"

	//
	// The attribution model of the conversion type is uneditable.
	//
	ConversionTrackingErrorReasonATTRIBUTION_MODEL_IS_UNEDITABLE ConversionTrackingErrorReason = "ATTRIBUTION_MODEL_IS_UNEDITABLE"

	//
	// The attribution model of the conversion type cannot be unknown.
	//
	ConversionTrackingErrorReasonATTRIBUTION_MODEL_CANNOT_BE_UNKNOWN ConversionTrackingErrorReason = "ATTRIBUTION_MODEL_CANNOT_BE_UNKNOWN"

	//
	// The attribution model cannot be set to DATA_DRIVEN because a data-driven model has never been
	// generated.
	//
	ConversionTrackingErrorReasonDATA_DRIVEN_MODEL_WAS_NEVER_GENERATED ConversionTrackingErrorReason = "DATA_DRIVEN_MODEL_WAS_NEVER_GENERATED"

	//
	// The attribution model cannot be set to DATA_DRIVEN because the data-driven model is expired.
	//
	ConversionTrackingErrorReasonDATA_DRIVEN_MODEL_IS_EXPIRED ConversionTrackingErrorReason = "DATA_DRIVEN_MODEL_IS_EXPIRED"

	//
	// The attribution model cannot be set to DATA_DRIVEN because the data-driven model is stale.
	//
	ConversionTrackingErrorReasonDATA_DRIVEN_MODEL_IS_STALE ConversionTrackingErrorReason = "DATA_DRIVEN_MODEL_IS_STALE"

	//
	// The attribution model cannot be set to DATA_DRIVEN because the data-driven model is
	// unavailable or the conversion type was newly added.
	//
	ConversionTrackingErrorReasonDATA_DRIVEN_MODEL_IS_UNKNOWN ConversionTrackingErrorReason = "DATA_DRIVEN_MODEL_IS_UNKNOWN"

	//
	// An attempt to access a conversion type failed because no conversion type with this ID
	// exists for this account.
	//
	ConversionTrackingErrorReasonCONVERSION_TYPE_NOT_FOUND ConversionTrackingErrorReason = "CONVERSION_TYPE_NOT_FOUND"

	//
	// The user attempted to change the click-through conversion (ctc) lookback window when it is
	// not editable.
	//
	ConversionTrackingErrorReasonCTC_LOOKBACK_WINDOW_IS_UNEDITABLE ConversionTrackingErrorReason = "CTC_LOOKBACK_WINDOW_IS_UNEDITABLE"

	//
	// An exception occurred in the domain layer during an attempt to process a
	// ConversionTypeOperation.
	//
	ConversionTrackingErrorReasonDOMAIN_EXCEPTION ConversionTrackingErrorReason = "DOMAIN_EXCEPTION"

	//
	// An attempt was made to set a counting type inconsistent with other properties.
	// Currently, AppConversion with appConversionType = DOWNLOAD and appPlatform = ANDROID_MARKET
	// cannot have a countingType of MANY_PER_CLICK
	//
	ConversionTrackingErrorReasonINCONSISTENT_COUNTING_TYPE ConversionTrackingErrorReason = "INCONSISTENT_COUNTING_TYPE"

	//
	// The user specified two identical app ids when attempting to create or modify a
	// conversion type.
	//
	ConversionTrackingErrorReasonDUPLICATE_APP_ID ConversionTrackingErrorReason = "DUPLICATE_APP_ID"

	//
	// The user is trying to enter a double bidding conflict. A double bidding conflict is when 2
	// conversion types both measure downloads for the same app ID.
	//
	ConversionTrackingErrorReasonTWO_CONVERSION_TYPES_BIDDING_ON_SAME_APP_DOWNLOAD ConversionTrackingErrorReason = "TWO_CONVERSION_TYPES_BIDDING_ON_SAME_APP_DOWNLOAD"

	//
	// The user is trying to enter a double bidding conflict with the global type. The conversion
	// type being created/editied and the global type (or forked global download type) are both
	// measuring downloads for the same app ID.
	//
	ConversionTrackingErrorReasonCONVERSION_TYPE_BIDDING_ON_SAME_APP_DOWNLOAD_AS_GLOBAL_TYPE ConversionTrackingErrorReason = "CONVERSION_TYPE_BIDDING_ON_SAME_APP_DOWNLOAD_AS_GLOBAL_TYPE"

	//
	// The user specified two identical names when attempting to create or rename multiple
	// conversion types.
	//
	ConversionTrackingErrorReasonDUPLICATE_NAME ConversionTrackingErrorReason = "DUPLICATE_NAME"

	//
	// An error occurred while the server was sending the email.
	//
	ConversionTrackingErrorReasonEMAIL_FAILED ConversionTrackingErrorReason = "EMAIL_FAILED"

	//
	// The maximum number of active conversion types for this account has been exceeded.
	//
	ConversionTrackingErrorReasonEXCEEDED_CONVERSION_TYPE_LIMIT ConversionTrackingErrorReason = "EXCEEDED_CONVERSION_TYPE_LIMIT"

	//
	// The user requested to modify an existing conversion type, but did not supply an ID.
	//
	ConversionTrackingErrorReasonID_IS_NULL ConversionTrackingErrorReason = "ID_IS_NULL"

	//
	// App ids must adhere to valid Java package naming requirements.
	//
	ConversionTrackingErrorReasonINVALID_APP_ID ConversionTrackingErrorReason = "INVALID_APP_ID"

	//
	// App id can not be set to forked system-defined Android download conversion type.
	//
	ConversionTrackingErrorReasonCANNOT_SET_APP_ID ConversionTrackingErrorReason = "CANNOT_SET_APP_ID"

	//
	// The user attempted to set category which is not applicable to provided conversion type.
	//
	ConversionTrackingErrorReasonINVALID_CATEGORY ConversionTrackingErrorReason = "INVALID_CATEGORY"

	//
	// The user entered an invalid background color. The background color must be a valid
	// HTML hex color code, such as "99ccff".
	//
	ConversionTrackingErrorReasonINVALID_COLOR ConversionTrackingErrorReason = "INVALID_COLOR"

	//
	// The date range specified in the stats selector is invalid.
	//
	ConversionTrackingErrorReasonINVALID_DATE_RANGE ConversionTrackingErrorReason = "INVALID_DATE_RANGE"

	//
	// The email address of the sender or the recipient of a snippet email was invalid.
	//
	ConversionTrackingErrorReasonINVALID_EMAIL_ADDRESS ConversionTrackingErrorReason = "INVALID_EMAIL_ADDRESS"

	//
	// When providing a global conversion type id to fork from in an ADD operation,
	// the global conversion type id is not acceptable (i.e.: we don't allow this global conversion
	// type to be forked from)
	//
	ConversionTrackingErrorReasonINVALID_ORIGINAL_CONVERSION_TYPE_ID ConversionTrackingErrorReason = "INVALID_ORIGINAL_CONVERSION_TYPE_ID"

	//
	// The AppPlatform and AppConversionType must be set at the same time. It is an error to set
	// just one or the other.
	//
	ConversionTrackingErrorReasonMUST_SET_APP_PLATFORM_AND_APP_CONVERSION_TYPE_TOGETHER ConversionTrackingErrorReason = "MUST_SET_APP_PLATFORM_AND_APP_CONVERSION_TYPE_TOGETHER"

	//
	// The user attempted to create a new conversion type, or to rename an existing conversion type,
	// whose new name matches one of the other conversion types for his account.
	//
	ConversionTrackingErrorReasonNAME_ALREADY_EXISTS ConversionTrackingErrorReason = "NAME_ALREADY_EXISTS"

	//
	// The user asked to send a notification email, but specified no recipients.
	//
	ConversionTrackingErrorReasonNO_RECIPIENTS ConversionTrackingErrorReason = "NO_RECIPIENTS"

	//
	// The requested conversion type has no snippet, and thus its snippet email cannot be sent.
	//
	ConversionTrackingErrorReasonNO_SNIPPET ConversionTrackingErrorReason = "NO_SNIPPET"

	//
	// The requested date range contains too many webpages to be processed.
	//
	ConversionTrackingErrorReasonTOO_MANY_WEBPAGES ConversionTrackingErrorReason = "TOO_MANY_WEBPAGES"

	//
	// An unknown sorting type was specified in the selector.
	//
	ConversionTrackingErrorReasonUNKNOWN_SORTING_TYPE ConversionTrackingErrorReason = "UNKNOWN_SORTING_TYPE"

	//
	// AppConversionType cannot be set to DOWNLOAD when AppPlatform is ITUNES.
	//
	ConversionTrackingErrorReasonUNSUPPORTED_APP_CONVERSION_TYPE ConversionTrackingErrorReason = "UNSUPPORTED_APP_CONVERSION_TYPE"
)

//
// Enumerates data driven model statuses.
//
type DataDrivenModelStatus string

const (

	//
	// The data driven model status is unknown.
	//
	DataDrivenModelStatusUNKNOWN DataDrivenModelStatus = "UNKNOWN"

	//
	// A data driven model is available.
	//
	DataDrivenModelStatusAVAILABLE DataDrivenModelStatus = "AVAILABLE"

	//
	// The data driven model is stale. It hasn't been updated for at least 7 days. It
	// is still being used, but will become expired if it does not get updated for 30
	// days.
	//
	DataDrivenModelStatusSTALE DataDrivenModelStatus = "STALE"

	//
	// The data driven model expired. It hasn't been updated for at least 30 days and
	// cannot be used. Most commonly this is because there haven't been the required
	// number of events in a recent 30-day period.
	//
	DataDrivenModelStatusEXPIRED DataDrivenModelStatus = "EXPIRED"

	//
	// A data driven model has never been generated. Most commonly this is because
	// there has never been the required number of events in any 30-day period.
	//
	DataDrivenModelStatusNEVER_GENERATED DataDrivenModelStatus = "NEVER_GENERATED"
)

//
// The reasons for the database error.
//
type DatabaseErrorReason string

const (

	//
	// A concurrency problem occurred as two threads were attempting to modify same object.
	//
	DatabaseErrorReasonCONCURRENT_MODIFICATION DatabaseErrorReason = "CONCURRENT_MODIFICATION"

	//
	// The permission was denied to access an object.
	//
	DatabaseErrorReasonPERMISSION_DENIED DatabaseErrorReason = "PERMISSION_DENIED"

	//
	// The user's access to an object has been prohibited.
	//
	DatabaseErrorReasonACCESS_PROHIBITED DatabaseErrorReason = "ACCESS_PROHIBITED"

	//
	// Requested campaign belongs to a product that is not supported by the api.
	//
	DatabaseErrorReasonCAMPAIGN_PRODUCT_NOT_SUPPORTED DatabaseErrorReason = "CAMPAIGN_PRODUCT_NOT_SUPPORTED"

	//
	// a duplicate key was detected upon insertion
	//
	DatabaseErrorReasonDUPLICATE_KEY DatabaseErrorReason = "DUPLICATE_KEY"

	//
	// a database error has occurred
	//
	DatabaseErrorReasonDATABASE_ERROR DatabaseErrorReason = "DATABASE_ERROR"

	//
	// an unknown error has occurred
	//
	DatabaseErrorReasonUNKNOWN DatabaseErrorReason = "UNKNOWN"
)

//
// The reasons for the target error.
//
type DateErrorReason string

const (

	//
	// Given field values do not correspond to a valid date.
	//
	DateErrorReasonINVALID_FIELD_VALUES_IN_DATE DateErrorReason = "INVALID_FIELD_VALUES_IN_DATE"

	//
	// Given field values do not correspond to a valid date time.
	//
	DateErrorReasonINVALID_FIELD_VALUES_IN_DATE_TIME DateErrorReason = "INVALID_FIELD_VALUES_IN_DATE_TIME"

	//
	// The string date's format should be yyyymmdd.
	//
	DateErrorReasonINVALID_STRING_DATE DateErrorReason = "INVALID_STRING_DATE"

	//
	// The string date range's format should be yyyymmdd yyyymmdd.
	//
	DateErrorReasonINVALID_STRING_DATE_RANGE DateErrorReason = "INVALID_STRING_DATE_RANGE"

	//
	// The string date time's format should be yyyymmdd hhmmss [tz].
	//
	DateErrorReasonINVALID_STRING_DATE_TIME DateErrorReason = "INVALID_STRING_DATE_TIME"

	//
	// Date is before allowed minimum.
	//
	DateErrorReasonEARLIER_THAN_MINIMUM_DATE DateErrorReason = "EARLIER_THAN_MINIMUM_DATE"

	//
	// Date is after allowed maximum.
	//
	DateErrorReasonLATER_THAN_MAXIMUM_DATE DateErrorReason = "LATER_THAN_MAXIMUM_DATE"

	//
	// Date range bounds are not in order.
	//
	DateErrorReasonDATE_RANGE_MINIMUM_DATE_LATER_THAN_MAXIMUM_DATE DateErrorReason = "DATE_RANGE_MINIMUM_DATE_LATER_THAN_MAXIMUM_DATE"

	//
	// Both dates in range are null.
	//
	DateErrorReasonDATE_RANGE_MINIMUM_AND_MAXIMUM_DATES_BOTH_NULL DateErrorReason = "DATE_RANGE_MINIMUM_AND_MAXIMUM_DATES_BOTH_NULL"
)

//
// The reasons for the validation error.
//
type DistinctErrorReason string

const (
	DistinctErrorReasonDUPLICATE_ELEMENT DistinctErrorReason = "DUPLICATE_ELEMENT"

	DistinctErrorReasonDUPLICATE_TYPE DistinctErrorReason = "DUPLICATE_TYPE"
)

//
// The reasons for the target error.
//
type IdErrorReason string

const (

	//
	// Id not found
	//
	IdErrorReasonNOT_FOUND IdErrorReason = "NOT_FOUND"
)

//
// The single reason for the internal API error.
//
type InternalApiErrorReason string

const (

	//
	// API encountered an unexpected internal error.
	//
	InternalApiErrorReasonUNEXPECTED_INTERNAL_API_ERROR InternalApiErrorReason = "UNEXPECTED_INTERNAL_API_ERROR"

	//
	// A temporary error occurred during the request. Please retry.
	//
	InternalApiErrorReasonTRANSIENT_ERROR InternalApiErrorReason = "TRANSIENT_ERROR"

	//
	// The cause of the error is not known or only defined in newer versions.
	//
	InternalApiErrorReasonUNKNOWN InternalApiErrorReason = "UNKNOWN"

	//
	// The API is currently unavailable for a planned downtime.
	//
	InternalApiErrorReasonDOWNTIME InternalApiErrorReason = "DOWNTIME"

	//
	// Mutate succeeded but server was unable to build response. Client should not retry mutate.
	//
	InternalApiErrorReasonERROR_GENERATING_RESPONSE InternalApiErrorReason = "ERROR_GENERATING_RESPONSE"
)

//
// The reasons for the validation error.
//
type NotEmptyErrorReason string

const (
	NotEmptyErrorReasonEMPTY_LIST NotEmptyErrorReason = "EMPTY_LIST"
)

//
// The reasons for the operation access error.
//
type OperationAccessDeniedReason string

const (

	//
	// Unauthorized invocation of a service's method (get, mutate, etc.)
	//
	OperationAccessDeniedReasonACTION_NOT_PERMITTED OperationAccessDeniedReason = "ACTION_NOT_PERMITTED"

	//
	// Unauthorized ADD operation in invoking a service's mutate method.
	//
	OperationAccessDeniedReasonADD_OPERATION_NOT_PERMITTED OperationAccessDeniedReason = "ADD_OPERATION_NOT_PERMITTED"

	//
	// Unauthorized REMOVE operation in invoking a service's mutate method.
	//
	OperationAccessDeniedReasonREMOVE_OPERATION_NOT_PERMITTED OperationAccessDeniedReason = "REMOVE_OPERATION_NOT_PERMITTED"

	//
	// Unauthorized SET operation in invoking a service's mutate method.
	//
	OperationAccessDeniedReasonSET_OPERATION_NOT_PERMITTED OperationAccessDeniedReason = "SET_OPERATION_NOT_PERMITTED"

	//
	// A mutate action is not allowed on this campaign, from this client.
	//
	OperationAccessDeniedReasonMUTATE_ACTION_NOT_PERMITTED_FOR_CLIENT OperationAccessDeniedReason = "MUTATE_ACTION_NOT_PERMITTED_FOR_CLIENT"

	//
	// This operation is not permitted on this campaign type
	//
	OperationAccessDeniedReasonOPERATION_NOT_PERMITTED_FOR_CAMPAIGN_TYPE OperationAccessDeniedReason = "OPERATION_NOT_PERMITTED_FOR_CAMPAIGN_TYPE"

	//
	// An ADD operation may not set status to REMOVED.
	//
	OperationAccessDeniedReasonADD_AS_REMOVED_NOT_PERMITTED OperationAccessDeniedReason = "ADD_AS_REMOVED_NOT_PERMITTED"

	//
	// This operation is not allowed because the campaign or adgroup is removed.
	//
	OperationAccessDeniedReasonOPERATION_NOT_PERMITTED_FOR_REMOVED_ENTITY OperationAccessDeniedReason = "OPERATION_NOT_PERMITTED_FOR_REMOVED_ENTITY"

	//
	// This operation is not permitted on this ad group type.
	//
	OperationAccessDeniedReasonOPERATION_NOT_PERMITTED_FOR_AD_GROUP_TYPE OperationAccessDeniedReason = "OPERATION_NOT_PERMITTED_FOR_AD_GROUP_TYPE"

	//
	// The reason the invoked method or operation is prohibited is not known
	// (the client may be of an older version than the server).
	//
	OperationAccessDeniedReasonUNKNOWN OperationAccessDeniedReason = "UNKNOWN"
)

//
// This represents an operator that may be presented to an adsapi service.
//
type Operator string

const (

	//
	// The ADD operator.
	//
	OperatorADD Operator = "ADD"

	//
	// The REMOVE operator.
	//
	OperatorREMOVE Operator = "REMOVE"

	//
	// The SET operator (used for updates).
	//
	OperatorSET Operator = "SET"
)

//
// The reasons for the validation error.
//
type OperatorErrorReason string

const (
	OperatorErrorReasonOPERATOR_NOT_SUPPORTED OperatorErrorReason = "OPERATOR_NOT_SUPPORTED"
)

//
// Defines the valid set of operators.
//
type PredicateOperator string

const (

	//
	// Checks if the field is equal to the given value.
	//
	// <p>This operator is used with integers, dates, booleans, strings,
	// enums, and sets.
	//
	PredicateOperatorEQUALS PredicateOperator = "EQUALS"

	//
	// Checks if the field does not equal the given value.
	//
	// <p>This operator is used with integers, booleans, strings, enums,
	// and sets.
	//
	PredicateOperatorNOT_EQUALS PredicateOperator = "NOT_EQUALS"

	//
	// Checks if the field is equal to one of the given values.
	//
	// <p>This operator accepts multiple operands and is used with
	// integers, booleans, strings, and enums.
	//
	PredicateOperatorIN PredicateOperator = "IN"

	//
	// Checks if the field does not equal any of the given values.
	//
	// <p>This operator accepts multiple operands and is used with
	// integers, booleans, strings, and enums.
	//
	PredicateOperatorNOT_IN PredicateOperator = "NOT_IN"

	//
	// Checks if the field is greater than the given value.
	//
	// <p>This operator is used with numbers and dates.
	//
	PredicateOperatorGREATER_THAN PredicateOperator = "GREATER_THAN"

	//
	// Checks if the field is greater or equal to the given value.
	//
	// <p>This operator is used with numbers and dates.
	//
	PredicateOperatorGREATER_THAN_EQUALS PredicateOperator = "GREATER_THAN_EQUALS"

	//
	// Checks if the field is less than the given value.
	//
	// <p>This operator is used with numbers and dates.
	//
	PredicateOperatorLESS_THAN PredicateOperator = "LESS_THAN"

	//
	// Checks if the field is less or equal to than the given value.
	//
	// <p>This operator is used with numbers and dates.
	//
	PredicateOperatorLESS_THAN_EQUALS PredicateOperator = "LESS_THAN_EQUALS"

	//
	// Checks if the field starts with the given value.
	//
	// <p>This operator is used with strings.
	//
	PredicateOperatorSTARTS_WITH PredicateOperator = "STARTS_WITH"

	//
	// Checks if the field starts with the given value, ignoring case.
	//
	// <p>This operator is used with strings.
	//
	PredicateOperatorSTARTS_WITH_IGNORE_CASE PredicateOperator = "STARTS_WITH_IGNORE_CASE"

	//
	// Checks if the field contains the given value as a substring.
	//
	// <p>This operator is used with strings.
	//
	PredicateOperatorCONTAINS PredicateOperator = "CONTAINS"

	//
	// Checks if the field contains the given value as a substring, ignoring
	// case.
	//
	// <p>This operator is used with strings.
	//
	PredicateOperatorCONTAINS_IGNORE_CASE PredicateOperator = "CONTAINS_IGNORE_CASE"

	//
	// Checks if the field does not contain the given value as a substring.
	//
	// <p>This operator is used with strings.
	//
	PredicateOperatorDOES_NOT_CONTAIN PredicateOperator = "DOES_NOT_CONTAIN"

	//
	// Checks if the field does not contain the given value as a substring,
	// ignoring case.
	//
	// <p>This operator is used with strings.
	//
	PredicateOperatorDOES_NOT_CONTAIN_IGNORE_CASE PredicateOperator = "DOES_NOT_CONTAIN_IGNORE_CASE"

	//
	// Checks if the field contains <em>any</em> of the given values.
	//
	// <p>This operator accepts multiple values and is used on sets of numbers
	// or strings.
	//
	PredicateOperatorCONTAINS_ANY PredicateOperator = "CONTAINS_ANY"

	//
	// Checks if the field contains <em>all</em> of the given values.
	//
	// <p>This operator accepts multiple values and is used on sets of numbers
	// or strings.
	//
	PredicateOperatorCONTAINS_ALL PredicateOperator = "CONTAINS_ALL"

	//
	// Checks if the field contains <em>none</em> of the given values.
	//
	// <p>This operator accepts multiple values and is used on sets of numbers
	// or strings.
	//
	PredicateOperatorCONTAINS_NONE PredicateOperator = "CONTAINS_NONE"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	PredicateOperatorUNKNOWN PredicateOperator = "UNKNOWN"
)

//
// The reason for the query error.
//
type QueryErrorReason string

const (

	//
	// Exception that happens when trying to parse a query that doesn't match the AWQL grammar.
	//
	QueryErrorReasonPARSING_FAILED QueryErrorReason = "PARSING_FAILED"

	//
	// The provided query is an empty string.
	//
	QueryErrorReasonMISSING_QUERY QueryErrorReason = "MISSING_QUERY"

	//
	// The query does not contain the required SELECT clause or it is not in the
	// correct location.
	//
	QueryErrorReasonMISSING_SELECT_CLAUSE QueryErrorReason = "MISSING_SELECT_CLAUSE"

	//
	// The query does not contain the required FROM clause or it is not in the correct location.
	//
	QueryErrorReasonMISSING_FROM_CLAUSE QueryErrorReason = "MISSING_FROM_CLAUSE"

	//
	// The SELECT clause could not be parsed.
	//
	QueryErrorReasonINVALID_SELECT_CLAUSE QueryErrorReason = "INVALID_SELECT_CLAUSE"

	//
	// The FROM clause could not be parsed.
	//
	QueryErrorReasonINVALID_FROM_CLAUSE QueryErrorReason = "INVALID_FROM_CLAUSE"

	//
	// The WHERE clause could not be parsed.
	//
	QueryErrorReasonINVALID_WHERE_CLAUSE QueryErrorReason = "INVALID_WHERE_CLAUSE"

	//
	// The ORDER BY clause could not be parsed.
	//
	QueryErrorReasonINVALID_ORDER_BY_CLAUSE QueryErrorReason = "INVALID_ORDER_BY_CLAUSE"

	//
	// The LIMIT clause could not be parsed.
	//
	QueryErrorReasonINVALID_LIMIT_CLAUSE QueryErrorReason = "INVALID_LIMIT_CLAUSE"

	//
	// The startIndex in the LIMIT clause does not contain a valid integer.
	//
	QueryErrorReasonINVALID_START_INDEX_IN_LIMIT_CLAUSE QueryErrorReason = "INVALID_START_INDEX_IN_LIMIT_CLAUSE"

	//
	// The pageSize in the LIMIT clause does not contain a valid integer.
	//
	QueryErrorReasonINVALID_PAGE_SIZE_IN_LIMIT_CLAUSE QueryErrorReason = "INVALID_PAGE_SIZE_IN_LIMIT_CLAUSE"

	//
	// The DURING clause could not be parsed.
	//
	QueryErrorReasonINVALID_DURING_CLAUSE QueryErrorReason = "INVALID_DURING_CLAUSE"

	//
	// The minimum date in the DURING clause is not a valid date in YYYYMMDD format.
	//
	QueryErrorReasonINVALID_MIN_DATE_IN_DURING_CLAUSE QueryErrorReason = "INVALID_MIN_DATE_IN_DURING_CLAUSE"

	//
	// The maximum date in the DURING clause is not a valid date in YYYYMMDD format.
	//
	QueryErrorReasonINVALID_MAX_DATE_IN_DURING_CLAUSE QueryErrorReason = "INVALID_MAX_DATE_IN_DURING_CLAUSE"

	//
	// The minimum date in the DURING is after the maximum date.
	//
	QueryErrorReasonMAX_LESS_THAN_MIN_IN_DURING_CLAUSE QueryErrorReason = "MAX_LESS_THAN_MIN_IN_DURING_CLAUSE"

	//
	// The query matched the grammar, but is invalid in some way such as using a service that
	// isn't supported.
	//
	QueryErrorReasonVALIDATION_FAILED QueryErrorReason = "VALIDATION_FAILED"
)

//
// Enums for all the reasons an error can be thrown to the user during
// billing quota checks.
//
type QuotaCheckErrorReason string

const (

	//
	// Customer passed in an invalid token in the header.
	//
	QuotaCheckErrorReasonINVALID_TOKEN_HEADER QuotaCheckErrorReason = "INVALID_TOKEN_HEADER"

	//
	// Customer is marked delinquent.
	//
	QuotaCheckErrorReasonACCOUNT_DELINQUENT QuotaCheckErrorReason = "ACCOUNT_DELINQUENT"

	//
	// Customer is a fraudulent.
	//
	QuotaCheckErrorReasonACCOUNT_INACCESSIBLE QuotaCheckErrorReason = "ACCOUNT_INACCESSIBLE"

	//
	// Inactive Account.
	//
	QuotaCheckErrorReasonACCOUNT_INACTIVE QuotaCheckErrorReason = "ACCOUNT_INACTIVE"

	//
	// Signup not complete
	//
	QuotaCheckErrorReasonINCOMPLETE_SIGNUP QuotaCheckErrorReason = "INCOMPLETE_SIGNUP"

	//
	// Developer token is not approved for production access, and the customer
	// is attempting to access a production account.
	//
	QuotaCheckErrorReasonDEVELOPER_TOKEN_NOT_APPROVED QuotaCheckErrorReason = "DEVELOPER_TOKEN_NOT_APPROVED"

	//
	// Terms and conditions are not signed.
	//
	QuotaCheckErrorReasonTERMS_AND_CONDITIONS_NOT_SIGNED QuotaCheckErrorReason = "TERMS_AND_CONDITIONS_NOT_SIGNED"

	//
	// Monthly budget quota reached.
	//
	QuotaCheckErrorReasonMONTHLY_BUDGET_REACHED QuotaCheckErrorReason = "MONTHLY_BUDGET_REACHED"

	//
	// Monthly budget quota exceeded.
	//
	QuotaCheckErrorReasonQUOTA_EXCEEDED QuotaCheckErrorReason = "QUOTA_EXCEEDED"
)

//
// The reasons for the target error.
//
type RangeErrorReason string

const (
	RangeErrorReasonTOO_LOW RangeErrorReason = "TOO_LOW"

	RangeErrorReasonTOO_HIGH RangeErrorReason = "TOO_HIGH"
)

//
// The reason for the rate exceeded error.
//
type RateExceededErrorReason string

const (

	//
	// Rate exceeded.
	//
	RateExceededErrorReasonRATE_EXCEEDED RateExceededErrorReason = "RATE_EXCEEDED"
)

//
// The reasons for the target error.
//
type ReadOnlyErrorReason string

const (
	ReadOnlyErrorReasonREAD_ONLY ReadOnlyErrorReason = "READ_ONLY"
)

//
// The reasons for the target error.
//
type RejectedErrorReason string

const (

	//
	// Unknown value encountered
	//
	RejectedErrorReasonUNKNOWN_VALUE RejectedErrorReason = "UNKNOWN_VALUE"
)

type RequestErrorReason string

const (

	//
	// Error reason is unknown.
	//
	RequestErrorReasonUNKNOWN RequestErrorReason = "UNKNOWN"

	//
	// Invalid input.
	//
	RequestErrorReasonINVALID_INPUT RequestErrorReason = "INVALID_INPUT"

	//
	// The api version in the request has been discontinued. Please update
	// to the new AdWords API version.
	//
	RequestErrorReasonUNSUPPORTED_VERSION RequestErrorReason = "UNSUPPORTED_VERSION"
)

//
// The reasons for the target error.
//
type RequiredErrorReason string

const (

	//
	// Missing required field.
	//
	RequiredErrorReasonREQUIRED RequiredErrorReason = "REQUIRED"
)

//
// The reasons for the target error.
//
type SelectorErrorReason string

const (

	//
	// The field name is not valid.
	//
	SelectorErrorReasonINVALID_FIELD_NAME SelectorErrorReason = "INVALID_FIELD_NAME"

	//
	// The list of fields is null or empty.
	//
	SelectorErrorReasonMISSING_FIELDS SelectorErrorReason = "MISSING_FIELDS"

	//
	// The list of predicates is null or empty.
	//
	SelectorErrorReasonMISSING_PREDICATES SelectorErrorReason = "MISSING_PREDICATES"

	//
	// Predicate operator does not support multiple values. Multiple values are
	// supported only for {@link Predicate.Operator#IN} and {@link Predicate.Operator#NOT_IN}.
	//
	SelectorErrorReasonOPERATOR_DOES_NOT_SUPPORT_MULTIPLE_VALUES SelectorErrorReason = "OPERATOR_DOES_NOT_SUPPORT_MULTIPLE_VALUES"

	//
	// The predicate enum value is not valid.
	//
	SelectorErrorReasonINVALID_PREDICATE_ENUM_VALUE SelectorErrorReason = "INVALID_PREDICATE_ENUM_VALUE"

	//
	// The predicate operator is empty.
	//
	SelectorErrorReasonMISSING_PREDICATE_OPERATOR SelectorErrorReason = "MISSING_PREDICATE_OPERATOR"

	//
	// The predicate values are empty.
	//
	SelectorErrorReasonMISSING_PREDICATE_VALUES SelectorErrorReason = "MISSING_PREDICATE_VALUES"

	//
	// The predicate field name is not valid.
	//
	SelectorErrorReasonINVALID_PREDICATE_FIELD_NAME SelectorErrorReason = "INVALID_PREDICATE_FIELD_NAME"

	//
	// The predicate operator is not valid.
	//
	SelectorErrorReasonINVALID_PREDICATE_OPERATOR SelectorErrorReason = "INVALID_PREDICATE_OPERATOR"

	//
	// Invalid selection of fields.
	//
	SelectorErrorReasonINVALID_FIELD_SELECTION SelectorErrorReason = "INVALID_FIELD_SELECTION"

	//
	// The predicate value is not valid.
	//
	SelectorErrorReasonINVALID_PREDICATE_VALUE SelectorErrorReason = "INVALID_PREDICATE_VALUE"

	//
	// The sort field name is not valid or the field is not sortable.
	//
	SelectorErrorReasonINVALID_SORT_FIELD_NAME SelectorErrorReason = "INVALID_SORT_FIELD_NAME"

	//
	// Standard error.
	//
	SelectorErrorReasonSELECTOR_ERROR SelectorErrorReason = "SELECTOR_ERROR"

	//
	// Filtering by date range is not supported.
	//
	SelectorErrorReasonFILTER_BY_DATE_RANGE_NOT_SUPPORTED SelectorErrorReason = "FILTER_BY_DATE_RANGE_NOT_SUPPORTED"

	//
	// Selector paging start index is too high.
	//
	SelectorErrorReasonSTART_INDEX_IS_TOO_HIGH SelectorErrorReason = "START_INDEX_IS_TOO_HIGH"

	//
	// The values list in a predicate was too long.
	//
	SelectorErrorReasonTOO_MANY_PREDICATE_VALUES SelectorErrorReason = "TOO_MANY_PREDICATE_VALUES"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	SelectorErrorReasonUNKNOWN_ERROR SelectorErrorReason = "UNKNOWN_ERROR"
)

//
// The reasons for Ad Scheduling errors.
//
type SizeLimitErrorReason string

const (

	//
	// The number of entries in the request exceeds the system limit.
	//
	SizeLimitErrorReasonREQUEST_SIZE_LIMIT_EXCEEDED SizeLimitErrorReason = "REQUEST_SIZE_LIMIT_EXCEEDED"

	//
	// The number of entries in the response exceeds the system limit.
	//
	SizeLimitErrorReasonRESPONSE_SIZE_LIMIT_EXCEEDED SizeLimitErrorReason = "RESPONSE_SIZE_LIMIT_EXCEEDED"

	//
	// The account is too large to load.
	//
	SizeLimitErrorReasonINTERNAL_STORAGE_ERROR SizeLimitErrorReason = "INTERNAL_STORAGE_ERROR"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	SizeLimitErrorReasonUNKNOWN SizeLimitErrorReason = "UNKNOWN"
)

//
// Possible orders of sorting.
//
type SortOrder string

const (
	SortOrderASCENDING SortOrder = "ASCENDING"

	SortOrderDESCENDING SortOrder = "DESCENDING"
)

//
// The reasons for the target error.
//
type StringFormatErrorReason string

const (
	StringFormatErrorReasonUNKNOWN StringFormatErrorReason = "UNKNOWN"

	//
	// The input string value contains disallowed characters.
	//
	StringFormatErrorReasonILLEGAL_CHARS StringFormatErrorReason = "ILLEGAL_CHARS"

	//
	// The input string value is invalid for the associated field.
	//
	StringFormatErrorReasonINVALID_FORMAT StringFormatErrorReason = "INVALID_FORMAT"
)

//
// The reasons for the target error.
//
type StringLengthErrorReason string

const (
	StringLengthErrorReasonTOO_SHORT StringLengthErrorReason = "TOO_SHORT"

	StringLengthErrorReasonTOO_LONG StringLengthErrorReason = "TOO_LONG"
)

type Get struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 get"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	ServiceSelector *Selector `xml:"serviceSelector,omitempty"`
}

type GetResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getResponse"`

	Rval *ConversionTrackerPage `xml:"rval,omitempty"`
}

type Mutate struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutate"`

	//
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint DistinctIds">Elements in this field must have distinct IDs for following {@link Operator}s : SET, REMOVE.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint SupportedOperators">The following {@link Operator}s are supported: ADD, SET.</span>
	//
	Operations []*ConversionTrackerOperation `xml:"operations,omitempty"`
}

type MutateResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutateResponse"`

	Rval *ConversionTrackerReturnValue `xml:"rval,omitempty"`
}

type Query struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 query"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Query string `xml:"query,omitempty"`
}

type QueryResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 queryResponse"`

	Rval *ConversionTrackerPage `xml:"rval,omitempty"`
}

type AdCallMetricsConversion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdCallMetricsConversion"`

	*ConversionTracker

	//
	// The phone-call duration (in seconds) after which a conversion should be reported for this
	// AdCallMetricsConversion.
	// <span class="constraint Selectable">This field can be selected using the value "PhoneCallDuration".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint InRange">This field must be between 0 and 10000, inclusive.</span>
	//
	PhoneCallDuration int64 `xml:"phoneCallDuration,omitempty"`
}

type AdWordsConversionTracker struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdWordsConversionTracker"`

	*ConversionTracker

	//
	// Tracking code to use for the conversion type.
	// <span class="constraint Selectable">This field can be selected using the value "TrackingCodeType".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	TrackingCodeType *AdWordsConversionTrackerTrackingCodeType `xml:"trackingCodeType,omitempty"`
}

type ApiError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ApiError"`

	//
	// The OGNL field path to identify cause of error.
	//
	FieldPath string `xml:"fieldPath,omitempty"`

	//
	// A parsed copy of the field path. For example, the field path "operations[1].operand"
	// corresponds to this list: {FieldPathElement(field = "operations", index = 1),
	// FieldPathElement(field = "operand", index = null)}.
	//
	FieldPathElements []*FieldPathElement `xml:"fieldPathElements,omitempty"`

	//
	// The data that caused the error.
	//
	Trigger string `xml:"trigger,omitempty"`

	//
	// A simple string representation of the error and reason.
	//
	ErrorString string `xml:"errorString,omitempty"`

	//
	// Indicates that this instance is a subtype of ApiError.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	ApiErrorType string `xml:"ApiError.Type,omitempty"`
}

type ApiException struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ApiException"`

	*ApplicationException

	//
	// List of errors.
	//
	Errors []*ApiError `xml:"errors,omitempty"`
}

type AppConversion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AppConversion"`

	*ConversionTracker

	//
	// App ID of the app conversion tracker. This field is required for certain
	// conversion types, in particular, Android app install (first open) and
	// Android app install (from Google Play).
	// <span class="constraint Selectable">This field can be selected using the value "AppId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	AppId string `xml:"appId,omitempty"`

	//
	// App platform of the app conversion tracker. This field defaults to NONE.
	// Once it is set to a value other than NONE it cannot be changed again. It must be
	// set at the same time as AppConversionType.
	// <span class="constraint Selectable">This field can be selected using the value "AppPlatform".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	AppPlatform *AppConversionAppPlatform `xml:"appPlatform,omitempty"`

	//
	// The generated snippet for this conversion tracker. This snippet is
	// auto-generated by the API, and will be ignored in mutate operands. This
	// field will always be returned for conversion trackers using snippets.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Snippet string `xml:"snippet,omitempty"`

	//
	// The type of AppConversion, which identifies a conversion as being either download or
	// in-app purchase. This field can only be set once and future reads will populate the type
	// appropriately. It is an error to change the value once it is set. This field must be set
	// at the same time as AppPlatform.
	//
	AppConversionType *AppConversionAppConversionType `xml:"appConversionType,omitempty"`

	//
	// The postback URL. When the conversion type specifies a postback url,
	// Google will send information about each conversion event to that url as they happen.
	// Details, including formatting requirements for this field:
	// https://developers.google.com/app-conversion-tracking/docs/app-install-feedback
	// <span class="constraint Selectable">This field can be selected using the value "AppPostbackUrl".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	AppPostbackUrl string `xml:"appPostbackUrl,omitempty"`
}

type AppPostbackUrlError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AppPostbackUrlError"`

	*ApiError

	Reason *AppPostbackUrlErrorReason `xml:"reason,omitempty"`
}

type ApplicationException struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ApplicationException"`

	//
	// Error message.
	//
	Message string `xml:"message,omitempty"`

	//
	// Indicates that this instance is a subtype of ApplicationException.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	ApplicationExceptionType string `xml:"ApplicationException.Type,omitempty"`
}

type AuthenticationError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AuthenticationError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AuthenticationErrorReason `xml:"reason,omitempty"`
}

type AuthorizationError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AuthorizationError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AuthorizationErrorReason `xml:"reason,omitempty"`
}

type ClientTermsError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ClientTermsError"`

	*ApiError

	Reason *ClientTermsErrorReason `xml:"reason,omitempty"`
}

type ConversionTrackerPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ConversionTrackerPage"`

	*NoStatsPage

	//
	// The result entries in this page.
	//
	Entries []*ConversionTracker `xml:"entries,omitempty"`
}

type ConversionTracker struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ConversionTracker"`

	//
	// ID of this conversion tracker, {@code null} when creating a new one.
	//
	// <p>There are some system-defined conversion trackers that are available
	// for all customers to use.  See {@link ConversionTrackerService#mutate} for
	// more information about how to modify these types.
	// <ul>
	// <li>179 - Calls from Ads</li>
	// <li>214 - Android Downloads</li>
	// <li>239 - Store Visits</li>
	// </ul>
	// <span class="constraint Selectable">This field can be selected using the value "Id".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: ADD.</span>
	//
	Id int64 `xml:"id,omitempty"`

	//
	// The ID of the original conversion type on which this ConversionType is based.
	// This is used to facilitate a connection between an existing shared conversion type
	// (e.g. Calls from ads) and an advertiser-specific conversion type. This may only be specified
	// for ADD operations, and can never be modified once a ConversionType is created.
	// <span class="constraint Selectable">This field can be selected using the value "OriginalConversionTypeId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: SET.</span>
	//
	OriginalConversionTypeId int64 `xml:"originalConversionTypeId,omitempty"`

	//
	// Name of this conversion tracker.
	// <span class="constraint Selectable">This field can be selected using the value "Name".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Name string `xml:"name,omitempty"`

	//
	// Status of this conversion tracker.
	// <span class="constraint Selectable">This field can be selected using the value "Status".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	Status *ConversionTrackerStatus `xml:"status,omitempty"`

	//
	// The category of conversion that is being tracked.
	// <span class="constraint Selectable">This field can be selected using the value "Category".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	Category *ConversionTrackerCategory `xml:"category,omitempty"`

	//
	// The event snippet that works with the global site tag to track actions that should be counted
	// as conversions. Returns an empty string if the conversion tracker does not use snippets.
	// <span class="constraint Selectable">This field can be selected using the value "GoogleEventSnippet".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	GoogleEventSnippet string `xml:"googleEventSnippet,omitempty"`

	//
	// The global site tag that adds visitors to your basic remarketing lists and sets new cookies on
	// your domain, which will store information about the ad click that brought a user to your
	// website. Returns an empty string if the conversion tracker does not use snippets.
	// <span class="constraint Selectable">This field can be selected using the value "GoogleGlobalSiteTag".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	GoogleGlobalSiteTag string `xml:"googleGlobalSiteTag,omitempty"`

	//
	// The status of the data-driven attribution model for the conversion type.
	// <span class="constraint Selectable">This field can be selected using the value "DataDrivenModelStatus".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DataDrivenModelStatus *DataDrivenModelStatus `xml:"dataDrivenModelStatus,omitempty"`

	//
	// The external customer ID of the conversion type owner, or 0 if this is a system-defined
	// conversion type. Only the conversion type owner may edit properties of the conversion type.
	// <span class="constraint Selectable">This field can be selected using the value "ConversionTypeOwnerCustomerId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ConversionTypeOwnerCustomerId int64 `xml:"conversionTypeOwnerCustomerId,omitempty"`

	//
	// Lookback window for view-through conversions in days. This is the length of
	// time in which a conversion without a click can be attributed to an
	// impression.
	// <span class="constraint Selectable">This field can be selected using the value "ViewthroughLookbackWindow".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint InRange">This field must be between 1 and 30, inclusive.</span>
	//
	ViewthroughLookbackWindow int32 `xml:"viewthroughLookbackWindow,omitempty"`

	//
	// The click-through conversion (ctc) lookback window is the maximum number of days between
	// the time a conversion event happens and the previous corresponding ad click.
	//
	// <p>Conversion events that occur more than this many days after the click are ignored.
	//
	// <p>This field is only editable for Adwords Conversions and Upload Conversions, but has a system
	// defined default for other types of conversions. The allowed range of values for this window
	// depends on the type of conversion and may expand, but 1-90 days is the currently allowed
	// range.
	// <span class="constraint Selectable">This field can be selected using the value "CtcLookbackWindow".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	CtcLookbackWindow int32 `xml:"ctcLookbackWindow,omitempty"`

	//
	// How to count events for this conversion tracker.
	// If countingType is MANY_PER_CLICK, then all conversion events are counted.
	// If countingType is ONE_PER_CLICK, then only the first conversion event of this type
	// following a given click will be counted.
	// More information is available at https://support.google.com/adwords/answer/3438531
	// <span class="constraint Selectable">This field can be selected using the value "CountingType".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	CountingType *ConversionDeduplicationMode `xml:"countingType,omitempty"`

	//
	// The value to use when the tag for this conversion tracker sends conversion events without
	// values. This value is applied on the server side, and is applicable to all ConversionTracker
	// subclasses.
	// <p>
	// See also the corresponding {@link ConversionTracker#defaultRevenueCurrencyCode}, and see
	// {@link ConversionTracker#alwaysUseDefaultRevenueValue} for details about when this value is
	// used.
	// <span class="constraint Selectable">This field can be selected using the value "DefaultRevenueValue".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint InRange">This field must be between 0 and 1000000000000, inclusive.</span>
	//
	DefaultRevenueValue float64 `xml:"defaultRevenueValue,omitempty"`

	//
	// The currency code to use when the tag for this conversion tracker sends conversion events
	// without currency codes. This code is applied on the server side, and is applicable to all
	// ConversionTracker subclasses. It must be a valid ISO4217 3-character code, such as USD.
	// <p>
	// This code is used if the code in the tag is not supplied or is unsupported, or if
	// {@link ConversionTracker#alwaysUseDefaultRevenueValue} is set to true. If this default code is
	// not set the currency code of the account is used as the default code.
	// <p>
	// Set the default code to XXX in order to specify that this conversion type does not have units
	// of a currency (that is, it is unitless). In this case no currency conversion will occur even if
	// a currency code is set in the tag.
	// <span class="constraint Selectable">This field can be selected using the value "DefaultRevenueCurrencyCode".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	DefaultRevenueCurrencyCode string `xml:"defaultRevenueCurrencyCode,omitempty"`

	//
	// Controls whether conversion event values and currency codes are taken from the tag snippet or
	// from {@link ConversionTracker#defaultRevenueValue} and
	// {@link ConversionTracker#defaultRevenueCurrencyCode}. If alwaysUseDefaultRevenueValue is true,
	// then conversion events will always use defaultRevenueValue and defaultRevenueCurrencyCode, even
	// if the tag has supplied a value and/or code when reporting the conversion event.  If
	// alwaysUseDefaultRevenueValue is false, then defaultRevenueValue and defaultRevenueCurrencyCode
	// are only used if the tag does not supply a value, or the tag's value is unparseable.
	// <span class="constraint Selectable">This field can be selected using the value "AlwaysUseDefaultRevenueValue".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	AlwaysUseDefaultRevenueValue bool `xml:"alwaysUseDefaultRevenueValue,omitempty"`

	//
	// Whether this conversion tracker should be excluded from the "Conversions" columns in reports.
	// <p>
	// If true, the conversion tracker will not be counted towards Conversions.
	// If false, it will be counted in Conversions. This is the default.</p>
	//
	// Either way, conversions will still be counted in the "AllConversions" columns in reports.
	// <span class="constraint Selectable">This field can be selected using the value "ExcludeFromBidding".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	ExcludeFromBidding bool `xml:"excludeFromBidding,omitempty"`

	//
	// Attribution models describing how to distribute credit for a particular conversion across
	// potentially many prior interactions. See https://support.google.com/adwords/answer/6259715 for
	// more information about attribution modeling in AdWords.
	// <span class="constraint Selectable">This field can be selected using the value "AttributionModelType".</span>
	//
	AttributionModelType *AttributionModelType `xml:"attributionModelType,omitempty"`

	//
	// The date of the most recent ad click that led to a conversion of this conversion type.
	//
	// <p>This date is in the <b>advertiser's defined time zone</b>.</p>
	// <span class="constraint Selectable">This field can be selected using the value "MostRecentConversionDate".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	MostRecentConversionDate string `xml:"mostRecentConversionDate,omitempty"`

	//
	// The last time a conversion tag for this conversion type successfully fired and was seen by
	// AdWords. This firing event may not have been the result of an attributable conversion
	// (ex: because the tag was fired from a browser that did not previously click an ad from the
	// appropriate advertiser).
	//
	// <p>This datetime is in <b>UTC</b>, not the advertiser's time zone.</p>
	// <span class="constraint Selectable">This field can be selected using the value "LastReceivedRequestTime".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	LastReceivedRequestTime string `xml:"lastReceivedRequestTime,omitempty"`

	//
	// Indicates that this instance is a subtype of ConversionTracker.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	ConversionTrackerType string `xml:"ConversionTracker.Type,omitempty"`
}

type ConversionTrackerOperation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ConversionTrackerOperation"`

	*Operation

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *ConversionTracker `xml:"operand,omitempty"`
}

type ConversionTrackerReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ConversionTrackerReturnValue"`

	*ListReturnValue

	Value []*ConversionTracker `xml:"value,omitempty"`
}

type ConversionTrackingError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ConversionTrackingError"`

	*ApiError

	Reason *ConversionTrackingErrorReason `xml:"reason,omitempty"`
}

type DatabaseError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DatabaseError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DatabaseErrorReason `xml:"reason,omitempty"`
}

type DateError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DateError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DateErrorReason `xml:"reason,omitempty"`
}

type DateRange struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DateRange"`

	//
	// the lower bound of this date range, inclusive.
	//
	Min string `xml:"min,omitempty"`

	//
	// the upper bound of this date range, inclusive.
	//
	Max string `xml:"max,omitempty"`
}

type DistinctError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DistinctError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DistinctErrorReason `xml:"reason,omitempty"`
}

type FieldPathElement struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 FieldPathElement"`

	//
	// The name of a field in lower camelcase. (e.g. "biddingStrategy")
	//
	Field string `xml:"field,omitempty"`

	//
	// For list fields, this is a 0-indexed position in the list. Null for non-list fields.
	//
	Index int32 `xml:"index,omitempty"`
}

type IdError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 IdError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *IdErrorReason `xml:"reason,omitempty"`
}

type InternalApiError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 InternalApiError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *InternalApiErrorReason `xml:"reason,omitempty"`
}

type ListReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ListReturnValue"`

	//
	// Indicates that this instance is a subtype of ListReturnValue.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	ListReturnValueType string `xml:"ListReturnValue.Type,omitempty"`
}

type NoStatsPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NoStatsPage"`

	*Page
}

type NotEmptyError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NotEmptyError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NotEmptyErrorReason `xml:"reason,omitempty"`
}

type Operation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Operation"`

	//
	// Operator.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operator *Operator `xml:"operator,omitempty"`

	//
	// Indicates that this instance is a subtype of Operation.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	OperationType string `xml:"Operation.Type,omitempty"`
}

type OperationAccessDenied struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 OperationAccessDenied"`

	*ApiError

	Reason *OperationAccessDeniedReason `xml:"reason,omitempty"`
}

type OperatorError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 OperatorError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *OperatorErrorReason `xml:"reason,omitempty"`
}

type OrderBy struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 OrderBy"`

	//
	// The field to sort the results on.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Field string `xml:"field,omitempty"`

	//
	// The order to sort the results on. The default sort order is {@link SortOrder#ASCENDING}.
	//
	SortOrder *SortOrder `xml:"sortOrder,omitempty"`
}

type Page struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Page"`

	//
	// Total number of entries in the result that this page is a part of.
	//
	TotalNumEntries int32 `xml:"totalNumEntries,omitempty"`

	//
	// Indicates that this instance is a subtype of Page.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	PageType string `xml:"Page.Type,omitempty"`
}

type Paging struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Paging"`

	//
	// Index of the first result to return in this page.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	StartIndex int32 `xml:"startIndex,omitempty"`

	//
	// Maximum number of results to return in this page. Set this to a reasonable value to limit
	// the number of results returned per page.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	NumberResults int32 `xml:"numberResults,omitempty"`
}

type Predicate struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Predicate"`

	//
	// The field by which to filter the returned data. Possible values are marked Filterable on
	// the entity's reference page. For example, for predicates for the
	// CampaignService {@link Selector selector}, refer to the filterable fields from the
	// {@link Campaign} reference page.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Field string `xml:"field,omitempty"`

	//
	// The operator to use for filtering the data returned.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operator *PredicateOperator `xml:"operator,omitempty"`

	//
	// The values by which to filter the field. The {@link Operator#CONTAINS_ALL},
	// {@link Operator#CONTAINS_ANY}, {@link Operator#CONTAINS_NONE}, {@link Operator#IN}
	// and {@link Operator#NOT_IN} take multiple values. All others take a single value.
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Values []string `xml:"values,omitempty"`
}

type QueryError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 QueryError"`

	*ApiError

	Reason *QueryErrorReason `xml:"reason,omitempty"`

	Message string `xml:"message,omitempty"`
}

type QuotaCheckError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 QuotaCheckError"`

	*ApiError

	Reason *QuotaCheckErrorReason `xml:"reason,omitempty"`
}

type RangeError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RangeError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *RangeErrorReason `xml:"reason,omitempty"`
}

type RateExceededError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RateExceededError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *RateExceededErrorReason `xml:"reason,omitempty"`

	//
	// Cause of the rate exceeded error.
	//
	RateName string `xml:"rateName,omitempty"`

	//
	// The scope of the rate (ACCOUNT/DEVELOPER).
	//
	RateScope string `xml:"rateScope,omitempty"`

	//
	// The amount of time (in seconds) the client should wait before retrying the request.
	//
	RetryAfterSeconds int32 `xml:"retryAfterSeconds,omitempty"`
}

type ReadOnlyError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ReadOnlyError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *ReadOnlyErrorReason `xml:"reason,omitempty"`
}

type RejectedError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RejectedError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *RejectedErrorReason `xml:"reason,omitempty"`
}

type RequestError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RequestError"`

	*ApiError

	Reason *RequestErrorReason `xml:"reason,omitempty"`
}

type RequiredError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RequiredError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *RequiredErrorReason `xml:"reason,omitempty"`
}

type Selector struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Selector"`

	//
	// List of fields to select.
	// <a href="/adwords/api/docs/appendix/selectorfields">Possible values</a>
	// are marked {@code Selectable} on the entity's reference page.
	// For example, for the {@code CampaignService} selector, refer to the
	// selectable fields from the {@link Campaign} reference page.
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Fields []string `xml:"fields,omitempty"`

	//
	// Specifies how an entity (eg. adgroup, campaign, criterion, ad) should be filtered.
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	//
	Predicates []*Predicate `xml:"predicates,omitempty"`

	//
	// Range of dates for which you want to include data. If this value is omitted,
	// results for all dates are returned.
	// <p class="note"><b>Note:</b> This field is only used by the report download
	// service. For all other services, it is ignored.</p>
	// <span class="constraint DateRangeWithinRange">This range must be contained within the range [19700101, 20380101].</span>
	//
	DateRange *DateRange `xml:"dateRange,omitempty"`

	//
	// The fields on which you want to sort, and the sort order. The order in the list is
	// significant: The first element in the list indicates the primary sort order, the next
	// specifies the secondary sort order and so on.
	//
	Ordering []*OrderBy `xml:"ordering,omitempty"`

	//
	// Pagination information.
	//
	Paging *Paging `xml:"paging,omitempty"`
}

type SelectorError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 SelectorError"`

	*ApiError

	//
	// The error reason represented by enum.
	//
	Reason *SelectorErrorReason `xml:"reason,omitempty"`
}

type SizeLimitError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 SizeLimitError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *SizeLimitErrorReason `xml:"reason,omitempty"`
}

type SoapHeader struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 SoapHeader"`

	//
	// The header identifies the customer id of the client of the AdWords manager, if an AdWords
	// manager is acting on behalf of their client or the customer id of the advertiser managing their
	// own account.
	//
	ClientCustomerId string `xml:"clientCustomerId,omitempty"`

	//
	// Developer token to identify that the person making the call has enough
	// quota.
	//
	DeveloperToken string `xml:"developerToken,omitempty"`

	//
	// UserAgent is used to track distribution of API client programs and
	// application usage. The client is responsible for putting in a meaningful
	// value for tracking purposes. To be clear this is not the same as an HTTP
	// user agent.
	//
	UserAgent string `xml:"userAgent,omitempty"`

	//
	// Used to validate the request without executing it.
	//
	ValidateOnly bool `xml:"validateOnly,omitempty"`

	//
	// If true, API will try to commit as many error free operations as possible and
	// report the other operations' errors.
	//
	// <p>Ignored for non-mutate calls.
	//
	PartialFailure bool `xml:"partialFailure,omitempty"`
}

type SoapResponseHeader struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 SoapResponseHeader"`

	//
	// Unique id that identifies this request. If developers have any support issues, sending us
	// this id will enable us to find their request more easily.
	//
	RequestId string `xml:"requestId,omitempty"`

	//
	// The name of the service being invoked.
	//
	ServiceName string `xml:"serviceName,omitempty"`

	//
	// The name of the method being invoked.
	//
	MethodName string `xml:"methodName,omitempty"`

	//
	// Number of operations performed for this SOAP request.
	//
	Operations int64 `xml:"operations,omitempty"`

	//
	// Elapsed time in milliseconds between the AdWords API receiving the request and sending the
	// response.
	//
	ResponseTime int64 `xml:"responseTime,omitempty"`
}

type StringFormatError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 StringFormatError"`

	*ApiError

	Reason *StringFormatErrorReason `xml:"reason,omitempty"`
}

type StringLengthError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 StringLengthError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *StringLengthErrorReason `xml:"reason,omitempty"`
}

type UploadCallConversion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 UploadCallConversion"`

	*ConversionTracker
}

type UploadConversion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 UploadConversion"`

	*ConversionTracker

	//
	// Whether this conversion is using an external attribution model.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: SET.</span>
	//
	IsExternallyAttributed bool `xml:"isExternallyAttributed,omitempty"`
}

type WebsiteCallMetricsConversion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 WebsiteCallMetricsConversion"`

	*ConversionTracker

	//
	// <span class="constraint Selectable">This field can be selected using the value "WebsitePhoneCallDuration".</span>
	// <span class="constraint InRange">This field must be between 0 and 10000, inclusive.</span>
	//
	PhoneCallDuration int64 `xml:"phoneCallDuration,omitempty"`
}

type ConversionTrackerServiceInterface struct {
	client *SOAPClient
}

func NewConversionTrackerServiceInterface(url string, tls bool, auth *BasicAuth) *ConversionTrackerServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &ConversionTrackerServiceInterface{
		client: client,
	}
}

func NewConversionTrackerServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *ConversionTrackerServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &ConversionTrackerServiceInterface{
		client: client,
	}
}

func (service *ConversionTrackerServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *ConversionTrackerServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Returns a list of the conversion trackers that match the selector. The
   actual objects contained in the page's list of entries will be specific
   subclasses of the abstract {@link ConversionTracker} class.

   @param serviceSelector The selector specifying the
   {@link ConversionTracker}s to return.
   @return List of conversion trackers specified by the selector.
   @throws com.google.ads.api.services.common.error.ApiException if problems
   occurred while retrieving results.
*/
func (service *ConversionTrackerServiceInterface) Get(request *Get) (*GetResponse, error) {
	response := new(GetResponse)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Applies the list of mutate operations such as adding or updating conversion trackers.
   <p class="note"><b>Note:</b> {@link ConversionTrackerOperation} does not support the
   <code>REMOVE</code> operator. In order to 'disable' a conversion type, send a
   <code>SET</code> operation for the conversion tracker with the <code>status</code>
   property set to <code>DISABLED</code></p>

   <p>You can mutate any ConversionTracker that belongs to your account. You may not
   mutate a ConversionTracker that belongs to some other account. You may not directly
   mutate a system-defined ConversionTracker, but you can create a mutable copy of it
   in your account by sending a mutate request with an ADD operation specifying
   an originalConversionTypeId matching a system-defined conversion tracker's ID. That new
   ADDed ConversionTracker will inherit the statistics and properties
   of the system-defined type, but will be editable as usual.</p>

   @param operations A list of mutate operations to perform.
   @return The list of the conversion trackers as they appear after mutation,
   in the same order as they appeared in the list of operations.
   @throws com.google.ads.api.services.common.error.ApiException if problems
   occurred while updating the data.
*/
func (service *ConversionTrackerServiceInterface) Mutate(request *Mutate) (*MutateResponse, error) {
	response := new(MutateResponse)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Returns a list of conversion trackers that match the query.

   @param query The SQL-like AWQL query string.
   @return A list of conversion trackers.
   @throws ApiException if problems occur while parsing the query or fetching conversion trackers.
*/
func (service *ConversionTrackerServiceInterface) Query(request *Query) (*QueryResponse, error) {
	response := new(QueryResponse)
	err := service.client.Call("", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  *SOAPHeader
	Body    SOAPBody
}

type SOAPHeader struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`

	Items []interface{} `xml:",omitempty"`
}

type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

const (
	// Predefined WSS namespaces to be used in
	WssNsWSSE string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
	WssNsWSU  string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd"
	WssNsType string = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordText"
)

type WSSSecurityHeader struct {
	XMLName   xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ wsse:Security"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	MustUnderstand string `xml:"mustUnderstand,attr,omitempty"`

	Token *WSSUsernameToken `xml:",omitempty"`
}

type WSSUsernameToken struct {
	XMLName   xml.Name `xml:"wsse:UsernameToken"`
	XmlNSWsu  string   `xml:"xmlns:wsu,attr"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	Id string `xml:"wsu:Id,attr,omitempty"`

	Username *WSSUsername `xml:",omitempty"`
	Password *WSSPassword `xml:",omitempty"`
}

type WSSUsername struct {
	XMLName   xml.Name `xml:"wsse:Username"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`

	Data string `xml:",chardata"`
}

type WSSPassword struct {
	XMLName   xml.Name `xml:"wsse:Password"`
	XmlNSWsse string   `xml:"xmlns:wsse,attr"`
	XmlNSType string   `xml:"Type,attr"`

	Data string `xml:",chardata"`
}

type BasicAuth struct {
	Login    string
	Password string
}

type SOAPClient struct {
	url     string
	tlsCfg  *tls.Config
	auth    *BasicAuth
	headers []interface{}
}

// **********
// Accepted solution from http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
// Author: Icza - http://stackoverflow.com/users/1705598/icza

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrc(n int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// **********

func NewWSSSecurityHeader(user, pass, mustUnderstand string) *WSSSecurityHeader {
	hdr := &WSSSecurityHeader{XmlNSWsse: WssNsWSSE, MustUnderstand: mustUnderstand}
	hdr.Token = &WSSUsernameToken{XmlNSWsu: WssNsWSU, XmlNSWsse: WssNsWSSE, Id: "UsernameToken-" + randStringBytesMaskImprSrc(9)}
	hdr.Token.Username = &WSSUsername{XmlNSWsse: WssNsWSSE, Data: user}
	hdr.Token.Password = &WSSPassword{XmlNSWsse: WssNsWSSE, XmlNSType: WssNsType, Data: pass}
	return hdr
}

func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &SOAPFault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}

func (f *SOAPFault) Error() string {
	return f.String
}

func NewSOAPClient(url string, insecureSkipVerify bool, auth *BasicAuth) *SOAPClient {
	tlsCfg := &tls.Config{
		InsecureSkipVerify: insecureSkipVerify,
	}
	return NewSOAPClientWithTLSConfig(url, tlsCfg, auth)
}

func NewSOAPClientWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *SOAPClient {
	return &SOAPClient{
		url:    url,
		tlsCfg: tlsCfg,
		auth:   auth,
	}
}

func (s *SOAPClient) AddHeader(header interface{}) {
	s.headers = append(s.headers, header)
}

func (s *SOAPClient) Call(soapAction string, request, response interface{}) error {
	envelope := SOAPEnvelope{}

	if s.headers != nil && len(s.headers) > 0 {
		soapHeader := &SOAPHeader{Items: make([]interface{}, len(s.headers))}
		copy(soapHeader.Items, s.headers)
		envelope.Header = soapHeader
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	log.Println(buffer.String())

	req, err := http.NewRequest("POST", s.url, buffer)
	if err != nil {
		return err
	}
	if s.auth != nil {
		req.SetBasicAuth(s.auth.Login, s.auth.Password)
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	req.Header.Add("SOAPAction", soapAction)

	req.Header.Set("User-Agent", "gowsdl/0.1")
	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: s.tlsCfg,
		Dial:            dialTimeout,
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawbody) == 0 {
		log.Println("empty response")
		return nil
	}

	log.Println(string(rawbody))
	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}
