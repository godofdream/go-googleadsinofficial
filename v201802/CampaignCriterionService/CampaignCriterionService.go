package CampaignCriterionService

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
// The reasons for the AdX error.
//
type AdxErrorReason string

const (

	//
	// Attempt to use non-AdX feature by AdX customer.
	//
	AdxErrorReasonUNSUPPORTED_FEATURE AdxErrorReason = "UNSUPPORTED_FEATURE"
)

type AgeRangeAgeRangeType string

const (
	AgeRangeAgeRangeTypeAGE_RANGE_18_24 AgeRangeAgeRangeType = "AGE_RANGE_18_24"

	AgeRangeAgeRangeTypeAGE_RANGE_25_34 AgeRangeAgeRangeType = "AGE_RANGE_25_34"

	AgeRangeAgeRangeTypeAGE_RANGE_35_44 AgeRangeAgeRangeType = "AGE_RANGE_35_44"

	AgeRangeAgeRangeTypeAGE_RANGE_45_54 AgeRangeAgeRangeType = "AGE_RANGE_45_54"

	AgeRangeAgeRangeTypeAGE_RANGE_55_64 AgeRangeAgeRangeType = "AGE_RANGE_55_64"

	AgeRangeAgeRangeTypeAGE_RANGE_65_UP AgeRangeAgeRangeType = "AGE_RANGE_65_UP"

	AgeRangeAgeRangeTypeAGE_RANGE_UNDETERMINED AgeRangeAgeRangeType = "AGE_RANGE_UNDETERMINED"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	AgeRangeAgeRangeTypeUNKNOWN AgeRangeAgeRangeType = "UNKNOWN"
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
// The status of the campaign criteria.
//
type CampaignCriterionCampaignCriterionStatus string

const (
	CampaignCriterionCampaignCriterionStatusACTIVE CampaignCriterionCampaignCriterionStatus = "ACTIVE"

	CampaignCriterionCampaignCriterionStatusREMOVED CampaignCriterionCampaignCriterionStatus = "REMOVED"

	CampaignCriterionCampaignCriterionStatusPAUSED CampaignCriterionCampaignCriterionStatus = "PAUSED"
)

//
// The reasons for the target error.
//
type CampaignCriterionErrorReason string

const (

	//
	// Concrete type of criterion (keyword v.s. placement) is required for
	// ADD and SET operations.
	//
	CampaignCriterionErrorReasonCONCRETE_TYPE_REQUIRED CampaignCriterionErrorReason = "CONCRETE_TYPE_REQUIRED"

	//
	// Invalid placement URL.
	//
	CampaignCriterionErrorReasonINVALID_PLACEMENT_URL CampaignCriterionErrorReason = "INVALID_PLACEMENT_URL"

	//
	// Criteria type can not be excluded for the campaign by the customer.
	// like AOL account type cannot target site type criteria
	//
	CampaignCriterionErrorReasonCANNOT_EXCLUDE_CRITERIA_TYPE CampaignCriterionErrorReason = "CANNOT_EXCLUDE_CRITERIA_TYPE"

	//
	// Cannot set the campaign criterion status for this criteria type.
	//
	CampaignCriterionErrorReasonCANNOT_SET_STATUS_FOR_CRITERIA_TYPE CampaignCriterionErrorReason = "CANNOT_SET_STATUS_FOR_CRITERIA_TYPE"

	//
	// Cannot set the campaign criterion status for an excluded criteria.
	//
	CampaignCriterionErrorReasonCANNOT_SET_STATUS_FOR_EXCLUDED_CRITERIA CampaignCriterionErrorReason = "CANNOT_SET_STATUS_FOR_EXCLUDED_CRITERIA"

	//
	// Cannot target and exclude the same criterion.
	//
	CampaignCriterionErrorReasonCANNOT_TARGET_AND_EXCLUDE CampaignCriterionErrorReason = "CANNOT_TARGET_AND_EXCLUDE"

	//
	// The #mutate operation contained too many operations.
	//
	CampaignCriterionErrorReasonTOO_MANY_OPERATIONS CampaignCriterionErrorReason = "TOO_MANY_OPERATIONS"

	//
	// This operator cannot be applied to a criterion of this type.
	//
	CampaignCriterionErrorReasonOPERATOR_NOT_SUPPORTED_FOR_CRITERION_TYPE CampaignCriterionErrorReason = "OPERATOR_NOT_SUPPORTED_FOR_CRITERION_TYPE"

	//
	// The Shopping campaign sales country is not supported for ProductSalesChannel targeting.
	//
	CampaignCriterionErrorReasonSHOPPING_CAMPAIGN_SALES_COUNTRY_NOT_SUPPORTED_FOR_SALES_CHANNEL CampaignCriterionErrorReason = "SHOPPING_CAMPAIGN_SALES_COUNTRY_NOT_SUPPORTED_FOR_SALES_CHANNEL"

	CampaignCriterionErrorReasonUNKNOWN CampaignCriterionErrorReason = "UNKNOWN"

	//
	// The existing field can't be updated with ADD operation. It can be updated with
	// SET operation only.
	//
	CampaignCriterionErrorReasonCANNOT_ADD_EXISTING_FIELD CampaignCriterionErrorReason = "CANNOT_ADD_EXISTING_FIELD"
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
// The types of constant operands.
//
type ConstantOperandConstantType string

const (

	//
	// Boolean constant type. booleanValue should be set for this type.
	//
	ConstantOperandConstantTypeBOOLEAN ConstantOperandConstantType = "BOOLEAN"

	//
	// Double constant type. doubleValue should be set for this type.
	//
	ConstantOperandConstantTypeDOUBLE ConstantOperandConstantType = "DOUBLE"

	//
	// Long constant type. longValue should be set for this type.
	//
	ConstantOperandConstantTypeLONG ConstantOperandConstantType = "LONG"

	//
	// String constant type. stringValue should be set for this type.
	//
	ConstantOperandConstantTypeSTRING ConstantOperandConstantType = "STRING"
)

//
// The units of constant operands, if applicable.
//
type ConstantOperandUnit string

const (

	//
	// Meters.
	//
	ConstantOperandUnitMETERS ConstantOperandUnit = "METERS"

	//
	// Miles.
	//
	ConstantOperandUnitMILES ConstantOperandUnit = "MILES"

	ConstantOperandUnitNONE ConstantOperandUnit = "NONE"
)

//
// Content label type.
//
type ContentLabelType string

const (

	//
	// Sexually suggestive content
	//
	ContentLabelTypeADULTISH ContentLabelType = "ADULTISH"

	//
	// Error pages
	//
	// <p class="note"><b>Note:</b> Starting with v201710, this label is deprecated and can only be
	// removed from campaigns - it can no longer be added. A future release will remove this label
	// entirely.
	//
	ContentLabelTypeAFE ContentLabelType = "AFE"

	//
	// Below the fold placements
	//
	ContentLabelTypeBELOW_THE_FOLD ContentLabelType = "BELOW_THE_FOLD"

	//
	// Military & international conflict
	//
	// <p class="note"><b>Note:</b> Starting with v201710, this label is deprecated and can only be
	// removed from campaigns - it can no longer be added. A future release will remove this label
	// entirely. Please use the {@code TRAGEDY} label instead of this one going forward.
	//
	ContentLabelTypeCONFLICT ContentLabelType = "CONFLICT"

	//
	// Parked domains
	//
	ContentLabelTypeDP ContentLabelType = "DP"

	//
	// Embedded video
	//
	ContentLabelTypeEMBEDDED_VIDEO ContentLabelType = "EMBEDDED_VIDEO"

	//
	// Games
	//
	ContentLabelTypeGAMES ContentLabelType = "GAMES"

	//
	// Sensational & shocking
	//
	ContentLabelTypeJUVENILE ContentLabelType = "JUVENILE"

	//
	// Profanity & rough language
	//
	ContentLabelTypePROFANITY ContentLabelType = "PROFANITY"

	//
	// Forums
	//
	// <p class="note"><b>Note:</b> Starting with v201710, this label is deprecated and can only be
	// removed from campaigns - it can no longer be added. A future release will remove this label
	// entirely.
	//
	ContentLabelTypeUGC_FORUMS ContentLabelType = "UGC_FORUMS"

	//
	// Image-sharing pages
	//
	// <p class="note"><b>Note:</b> Starting with v201710, this label is deprecated and can only be
	// removed from campaigns - it can no longer be added. A future release will remove this label
	// entirely.
	//
	ContentLabelTypeUGC_IMAGES ContentLabelType = "UGC_IMAGES"

	//
	// Social networks
	//
	// <p class="note"><b>Note:</b> Starting with v201710, this label is deprecated and can only be
	// removed from campaigns - it can no longer be added. A future release will remove this label
	// entirely.
	//
	ContentLabelTypeUGC_SOCIAL ContentLabelType = "UGC_SOCIAL"

	//
	// Video-sharing pages
	//
	// <p class="note"><b>Note:</b> Starting with v201710, this label is deprecated and can only be
	// removed from campaigns - it can no longer be added. A future release will remove this label
	// entirely.
	//
	ContentLabelTypeUGC_VIDEOS ContentLabelType = "UGC_VIDEOS"

	//
	// Crime, police & emergency
	//
	// <p class="note"><b>Note:</b> Starting with v201710, this label is deprecated and can only be
	// removed from campaigns - it can no longer be added. A future release will remove this label
	// entirely. Please use the {@code TRAGEDY} label instead of this one going forward.
	//
	ContentLabelTypeSIRENS ContentLabelType = "SIRENS"

	//
	// Tragedy & conflict
	//
	ContentLabelTypeTRAGEDY ContentLabelType = "TRAGEDY"

	//
	// Video
	//
	ContentLabelTypeVIDEO ContentLabelType = "VIDEO"

	//
	// Content rating: G
	//
	ContentLabelTypeVIDEO_RATING_DV_G ContentLabelType = "VIDEO_RATING_DV_G"

	//
	// Content rating: PG
	//
	ContentLabelTypeVIDEO_RATING_DV_PG ContentLabelType = "VIDEO_RATING_DV_PG"

	//
	// Content rating: T
	//
	ContentLabelTypeVIDEO_RATING_DV_T ContentLabelType = "VIDEO_RATING_DV_T"

	//
	// Content rating: MA
	//
	ContentLabelTypeVIDEO_RATING_DV_MA ContentLabelType = "VIDEO_RATING_DV_MA"

	//
	// Content rating: not yet rated
	//
	ContentLabelTypeVIDEO_NOT_YET_RATED ContentLabelType = "VIDEO_NOT_YET_RATED"

	//
	// Live streaming video
	//
	ContentLabelTypeLIVE_STREAMING_VIDEO ContentLabelType = "LIVE_STREAMING_VIDEO"

	//
	// Allowed gambling content.
	//
	// <p class="note"><b>Note:</b> Starting with v201710, this label is deprecated and can only be
	// removed from campaigns - it can no longer be added. A future release will remove this label
	// entirely.
	//
	ContentLabelTypeALLOWED_GAMBLING_CONTENT ContentLabelType = "ALLOWED_GAMBLING_CONTENT"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ContentLabelTypeUNKNOWN ContentLabelType = "UNKNOWN"
)

//
// The types of criteria.
//
type CriterionType string

const (

	//
	// Content label for exclusion.
	//
	CriterionTypeCONTENT_LABEL CriterionType = "CONTENT_LABEL"

	//
	// Keyword. e.g. 'mars cruise'
	//
	CriterionTypeKEYWORD CriterionType = "KEYWORD"

	//
	// Placement. aka Website. e.g. 'www.flowers4sale.com'
	//
	CriterionTypePLACEMENT CriterionType = "PLACEMENT"

	//
	// Vertical, e.g. 'category::Animals>Pets'  This is for vertical targeting on the content
	// network.
	//
	CriterionTypeVERTICAL CriterionType = "VERTICAL"

	//
	// User lists, are links to sets of users defined by the advertiser.
	//
	CriterionTypeUSER_LIST CriterionType = "USER_LIST"

	//
	// User interests, categories of interests the user is interested in.
	//
	CriterionTypeUSER_INTEREST CriterionType = "USER_INTEREST"

	//
	// Mobile applications to target.
	//
	CriterionTypeMOBILE_APPLICATION CriterionType = "MOBILE_APPLICATION"

	//
	// Mobile application categories to target.
	//
	CriterionTypeMOBILE_APP_CATEGORY CriterionType = "MOBILE_APP_CATEGORY"

	//
	// Product partition (product group) in a shopping campaign.
	//
	CriterionTypePRODUCT_PARTITION CriterionType = "PRODUCT_PARTITION"

	//
	// IP addresses to exclude.
	//
	CriterionTypeIP_BLOCK CriterionType = "IP_BLOCK"

	//
	// Webpages of an advertiser's website to target.
	//
	CriterionTypeWEBPAGE CriterionType = "WEBPAGE"

	//
	// Languages to target.
	//
	CriterionTypeLANGUAGE CriterionType = "LANGUAGE"

	//
	// Geographic regions to target.
	//
	CriterionTypeLOCATION CriterionType = "LOCATION"

	//
	// Age Range to exclude.
	//
	CriterionTypeAGE_RANGE CriterionType = "AGE_RANGE"

	//
	// Mobile carriers to target.
	//
	CriterionTypeCARRIER CriterionType = "CARRIER"

	//
	// Mobile operating system versions to target.
	//
	CriterionTypeOPERATING_SYSTEM_VERSION CriterionType = "OPERATING_SYSTEM_VERSION"

	//
	// Mobile devices to target.
	//
	CriterionTypeMOBILE_DEVICE CriterionType = "MOBILE_DEVICE"

	//
	// Gender to exclude.
	//
	CriterionTypeGENDER CriterionType = "GENDER"

	//
	// Parent to target and exclude.
	//
	CriterionTypePARENT CriterionType = "PARENT"

	//
	// Proximity (area within a radius) to target.
	//
	CriterionTypePROXIMITY CriterionType = "PROXIMITY"

	//
	// Platforms to target.
	//
	CriterionTypePLATFORM CriterionType = "PLATFORM"

	//
	// Representing preferred content bid modifier.
	//
	CriterionTypePREFERRED_CONTENT CriterionType = "PREFERRED_CONTENT"

	//
	// AdSchedule or specific days and time intervals to target.
	//
	CriterionTypeAD_SCHEDULE CriterionType = "AD_SCHEDULE"

	//
	// Targeting based on location groups.
	//
	CriterionTypeLOCATION_GROUPS CriterionType = "LOCATION_GROUPS"

	//
	// Scope of products. Contains a list of product dimensions, all of which a product has to match
	// to be included in the campaign.
	//
	CriterionTypePRODUCT_SCOPE CriterionType = "PRODUCT_SCOPE"

	//
	// YouTube video to target.
	//
	CriterionTypeYOUTUBE_VIDEO CriterionType = "YOUTUBE_VIDEO"

	//
	// YouTube channel to target.
	//
	CriterionTypeYOUTUBE_CHANNEL CriterionType = "YOUTUBE_CHANNEL"

	//
	// Enables advertisers to target paid apps.
	//
	CriterionTypeAPP_PAYMENT_MODEL CriterionType = "APP_PAYMENT_MODEL"

	//
	// Income range to target and exclude.
	//
	CriterionTypeINCOME_RANGE CriterionType = "INCOME_RANGE"

	//
	// Interaction type to bid modify.
	//
	CriterionTypeINTERACTION_TYPE CriterionType = "INTERACTION_TYPE"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	CriterionTypeUNKNOWN CriterionType = "UNKNOWN"
)

type CriterionErrorReason string

const (

	//
	// Concrete type of criterion is required for ADD and SET operations.
	//
	CriterionErrorReasonCONCRETE_TYPE_REQUIRED CriterionErrorReason = "CONCRETE_TYPE_REQUIRED"

	//
	// The category requested for exclusion is invalid.
	//
	CriterionErrorReasonINVALID_EXCLUDED_CATEGORY CriterionErrorReason = "INVALID_EXCLUDED_CATEGORY"

	//
	// Invalid keyword criteria text.
	//
	CriterionErrorReasonINVALID_KEYWORD_TEXT CriterionErrorReason = "INVALID_KEYWORD_TEXT"

	//
	// Keyword text should be less than 80 chars.
	//
	CriterionErrorReasonKEYWORD_TEXT_TOO_LONG CriterionErrorReason = "KEYWORD_TEXT_TOO_LONG"

	//
	// Keyword text has too many words.
	//
	CriterionErrorReasonKEYWORD_HAS_TOO_MANY_WORDS CriterionErrorReason = "KEYWORD_HAS_TOO_MANY_WORDS"

	//
	// Keyword text has invalid characters or symbols.
	//
	CriterionErrorReasonKEYWORD_HAS_INVALID_CHARS CriterionErrorReason = "KEYWORD_HAS_INVALID_CHARS"

	//
	// Invalid placement URL.
	//
	CriterionErrorReasonINVALID_PLACEMENT_URL CriterionErrorReason = "INVALID_PLACEMENT_URL"

	//
	// Invalid user list criterion.
	//
	CriterionErrorReasonINVALID_USER_LIST CriterionErrorReason = "INVALID_USER_LIST"

	//
	// Invalid user interest criterion.
	//
	CriterionErrorReasonINVALID_USER_INTEREST CriterionErrorReason = "INVALID_USER_INTEREST"

	//
	// Placement URL has wrong format.
	//
	CriterionErrorReasonINVALID_FORMAT_FOR_PLACEMENT_URL CriterionErrorReason = "INVALID_FORMAT_FOR_PLACEMENT_URL"

	//
	// Placement URL is too long.
	//
	CriterionErrorReasonPLACEMENT_URL_IS_TOO_LONG CriterionErrorReason = "PLACEMENT_URL_IS_TOO_LONG"

	//
	// Indicates the URL contains an illegal character.
	//
	CriterionErrorReasonPLACEMENT_URL_HAS_ILLEGAL_CHAR CriterionErrorReason = "PLACEMENT_URL_HAS_ILLEGAL_CHAR"

	//
	// Indicates the URL contains multiple comma separated URLs.
	//
	CriterionErrorReasonPLACEMENT_URL_HAS_MULTIPLE_SITES_IN_LINE CriterionErrorReason = "PLACEMENT_URL_HAS_MULTIPLE_SITES_IN_LINE"

	//
	// Indicates the domain is blacklisted.
	//
	CriterionErrorReasonPLACEMENT_IS_NOT_AVAILABLE_FOR_TARGETING_OR_EXCLUSION CriterionErrorReason = "PLACEMENT_IS_NOT_AVAILABLE_FOR_TARGETING_OR_EXCLUSION"

	//
	// Invalid vertical path.
	//
	CriterionErrorReasonINVALID_VERTICAL_PATH CriterionErrorReason = "INVALID_VERTICAL_PATH"

	//
	// Indicates the placement is a YouTube vertical channel, which is no longer supported.
	//
	CriterionErrorReasonYOUTUBE_VERTICAL_CHANNEL_DEPRECATED CriterionErrorReason = "YOUTUBE_VERTICAL_CHANNEL_DEPRECATED"

	//
	// Indicates the placement is a YouTube demographic channel, which is no longer supported.
	//
	CriterionErrorReasonYOUTUBE_DEMOGRAPHIC_CHANNEL_DEPRECATED CriterionErrorReason = "YOUTUBE_DEMOGRAPHIC_CHANNEL_DEPRECATED"

	//
	// YouTube urls are not supported in Placement criterion. Use YouTubeChannel and
	// YouTubeVideo criterion instead.
	//
	CriterionErrorReasonYOUTUBE_URL_UNSUPPORTED CriterionErrorReason = "YOUTUBE_URL_UNSUPPORTED"

	//
	// Criteria type can not be excluded by the customer,
	// like AOL account type cannot target site type criteria.
	//
	CriterionErrorReasonCANNOT_EXCLUDE_CRITERIA_TYPE CriterionErrorReason = "CANNOT_EXCLUDE_CRITERIA_TYPE"

	//
	// Criteria type can not be targeted.
	//
	CriterionErrorReasonCANNOT_ADD_CRITERIA_TYPE CriterionErrorReason = "CANNOT_ADD_CRITERIA_TYPE"

	//
	// Product filter in the product criteria has invalid characters.
	// Operand and the argument in the filter can not have "==" or "&+".
	//
	CriterionErrorReasonINVALID_PRODUCT_FILTER CriterionErrorReason = "INVALID_PRODUCT_FILTER"

	//
	// Product filter in the product criteria is translated to a string as
	// operand1==argument1&+operand2==argument2, maximum allowed length for
	// the string is 255 chars.
	//
	CriterionErrorReasonPRODUCT_FILTER_TOO_LONG CriterionErrorReason = "PRODUCT_FILTER_TOO_LONG"

	//
	// Not allowed to exclude similar user list.
	//
	CriterionErrorReasonCANNOT_EXCLUDE_SIMILAR_USER_LIST CriterionErrorReason = "CANNOT_EXCLUDE_SIMILAR_USER_LIST"

	//
	// Not allowed to target a closed user list.
	//
	CriterionErrorReasonCANNOT_ADD_CLOSED_USER_LIST CriterionErrorReason = "CANNOT_ADD_CLOSED_USER_LIST"

	//
	// Not allowed to add display only UserLists to search only campaigns.
	//
	CriterionErrorReasonCANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_ONLY_CAMPAIGNS CriterionErrorReason = "CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_ONLY_CAMPAIGNS"

	//
	// Not allowed to add display only UserLists to search plus campaigns.
	//
	CriterionErrorReasonCANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_CAMPAIGNS CriterionErrorReason = "CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_CAMPAIGNS"

	//
	// Not allowed to add display only UserLists to shopping campaigns.
	//
	CriterionErrorReasonCANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SHOPPING_CAMPAIGNS CriterionErrorReason = "CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SHOPPING_CAMPAIGNS"

	//
	// Not allowed to add User interests to search only campaigns.
	//
	CriterionErrorReasonCANNOT_ADD_USER_INTERESTS_TO_SEARCH_CAMPAIGNS CriterionErrorReason = "CANNOT_ADD_USER_INTERESTS_TO_SEARCH_CAMPAIGNS"

	//
	// Not allowed to set bids for this criterion type in search campaigns
	//
	CriterionErrorReasonCANNOT_SET_BIDS_ON_CRITERION_TYPE_IN_SEARCH_CAMPAIGNS CriterionErrorReason = "CANNOT_SET_BIDS_ON_CRITERION_TYPE_IN_SEARCH_CAMPAIGNS"

	//
	// Final URLs, URL Templates and CustomParameters cannot be set for the criterion
	// types of Gender, AgeRange, UserList, Placement, MobileApp, and MobileAppCategory
	// in search campaigns and shopping campaigns.
	//
	CriterionErrorReasonCANNOT_ADD_URLS_TO_CRITERION_TYPE_FOR_CAMPAIGN_TYPE CriterionErrorReason = "CANNOT_ADD_URLS_TO_CRITERION_TYPE_FOR_CAMPAIGN_TYPE"

	//
	// IP address is not valid.
	//
	CriterionErrorReasonINVALID_IP_ADDRESS CriterionErrorReason = "INVALID_IP_ADDRESS"

	//
	// IP format is not valid.
	//
	CriterionErrorReasonINVALID_IP_FORMAT CriterionErrorReason = "INVALID_IP_FORMAT"

	//
	// Mobile application is not valid.
	//
	CriterionErrorReasonINVALID_MOBILE_APP CriterionErrorReason = "INVALID_MOBILE_APP"

	//
	// Mobile application category is not valid.
	//
	CriterionErrorReasonINVALID_MOBILE_APP_CATEGORY CriterionErrorReason = "INVALID_MOBILE_APP_CATEGORY"

	//
	// The CriterionId does not exist or is of the incorrect type.
	//
	CriterionErrorReasonINVALID_CRITERION_ID CriterionErrorReason = "INVALID_CRITERION_ID"

	//
	// The Criterion is not allowed to be targeted.
	//
	CriterionErrorReasonCANNOT_TARGET_CRITERION CriterionErrorReason = "CANNOT_TARGET_CRITERION"

	//
	// The criterion is not allowed to be targeted as it is deprecated.
	//
	CriterionErrorReasonCANNOT_TARGET_OBSOLETE_CRITERION CriterionErrorReason = "CANNOT_TARGET_OBSOLETE_CRITERION"

	//
	// The CriterionId is not valid for the type.
	//
	CriterionErrorReasonCRITERION_ID_AND_TYPE_MISMATCH CriterionErrorReason = "CRITERION_ID_AND_TYPE_MISMATCH"

	//
	// Distance for the radius for the proximity criterion is invalid.
	//
	CriterionErrorReasonINVALID_PROXIMITY_RADIUS CriterionErrorReason = "INVALID_PROXIMITY_RADIUS"

	//
	// Units for the distance for the radius for the proximity criterion is invalid.
	//
	CriterionErrorReasonINVALID_PROXIMITY_RADIUS_UNITS CriterionErrorReason = "INVALID_PROXIMITY_RADIUS_UNITS"

	//
	// Street address is too short.
	//
	CriterionErrorReasonINVALID_STREETADDRESS_LENGTH CriterionErrorReason = "INVALID_STREETADDRESS_LENGTH"

	//
	// City name in the address is too short.
	//
	CriterionErrorReasonINVALID_CITYNAME_LENGTH CriterionErrorReason = "INVALID_CITYNAME_LENGTH"

	//
	// Region code in the address is too short.
	//
	CriterionErrorReasonINVALID_REGIONCODE_LENGTH CriterionErrorReason = "INVALID_REGIONCODE_LENGTH"

	//
	// Region name in the address is not valid.
	//
	CriterionErrorReasonINVALID_REGIONNAME_LENGTH CriterionErrorReason = "INVALID_REGIONNAME_LENGTH"

	//
	// Postal code in the address is not valid.
	//
	CriterionErrorReasonINVALID_POSTALCODE_LENGTH CriterionErrorReason = "INVALID_POSTALCODE_LENGTH"

	//
	// Country code in the address is not valid.
	//
	CriterionErrorReasonINVALID_COUNTRY_CODE CriterionErrorReason = "INVALID_COUNTRY_CODE"

	//
	// Latitude for the GeoPoint is not valid.
	//
	CriterionErrorReasonINVALID_LATITUDE CriterionErrorReason = "INVALID_LATITUDE"

	//
	// Longitude for the GeoPoint is not valid.
	//
	CriterionErrorReasonINVALID_LONGITUDE CriterionErrorReason = "INVALID_LONGITUDE"

	//
	// The Proximity input is not valid. Both address and geoPoint cannot be null.
	//
	CriterionErrorReasonPROXIMITY_GEOPOINT_AND_ADDRESS_BOTH_CANNOT_BE_NULL CriterionErrorReason = "PROXIMITY_GEOPOINT_AND_ADDRESS_BOTH_CANNOT_BE_NULL"

	//
	// The Proximity address cannot be geocoded to a valid lat/long.
	//
	CriterionErrorReasonINVALID_PROXIMITY_ADDRESS CriterionErrorReason = "INVALID_PROXIMITY_ADDRESS"

	//
	// User domain name is not valid.
	//
	CriterionErrorReasonINVALID_USER_DOMAIN_NAME CriterionErrorReason = "INVALID_USER_DOMAIN_NAME"

	//
	// Length of serialized criterion parameter exceeded size limit.
	//
	CriterionErrorReasonCRITERION_PARAMETER_TOO_LONG CriterionErrorReason = "CRITERION_PARAMETER_TOO_LONG"

	//
	// Time interval in the AdSchedule overlaps with another AdSchedule.
	//
	CriterionErrorReasonAD_SCHEDULE_TIME_INTERVALS_OVERLAP CriterionErrorReason = "AD_SCHEDULE_TIME_INTERVALS_OVERLAP"

	//
	// AdSchedule time interval cannot span multiple days.
	//
	CriterionErrorReasonAD_SCHEDULE_INTERVAL_CANNOT_SPAN_MULTIPLE_DAYS CriterionErrorReason = "AD_SCHEDULE_INTERVAL_CANNOT_SPAN_MULTIPLE_DAYS"

	//
	// AdSchedule time interval specified is invalid,
	// endTime cannot be earlier than startTime.
	//
	CriterionErrorReasonAD_SCHEDULE_INVALID_TIME_INTERVAL CriterionErrorReason = "AD_SCHEDULE_INVALID_TIME_INTERVAL"

	//
	// The number of AdSchedule entries in a day exceeds the limit.
	//
	CriterionErrorReasonAD_SCHEDULE_EXCEEDED_INTERVALS_PER_DAY_LIMIT CriterionErrorReason = "AD_SCHEDULE_EXCEEDED_INTERVALS_PER_DAY_LIMIT"

	//
	// CriteriaId does not match the interval of the AdSchedule specified.
	//
	CriterionErrorReasonAD_SCHEDULE_CRITERION_ID_MISMATCHING_FIELDS CriterionErrorReason = "AD_SCHEDULE_CRITERION_ID_MISMATCHING_FIELDS"

	//
	// Cannot set bid modifier for this criterion type.
	//
	CriterionErrorReasonCANNOT_BID_MODIFY_CRITERION_TYPE CriterionErrorReason = "CANNOT_BID_MODIFY_CRITERION_TYPE"

	//
	// Cannot bid modify criterion, since it is opted out of the campaign.
	//
	CriterionErrorReasonCANNOT_BID_MODIFY_CRITERION_CAMPAIGN_OPTED_OUT CriterionErrorReason = "CANNOT_BID_MODIFY_CRITERION_CAMPAIGN_OPTED_OUT"

	//
	// Cannot set bid modifier for a negative criterion.
	//
	CriterionErrorReasonCANNOT_BID_MODIFY_NEGATIVE_CRITERION CriterionErrorReason = "CANNOT_BID_MODIFY_NEGATIVE_CRITERION"

	//
	// Bid Modifier already exists. Use SET operation to update.
	//
	CriterionErrorReasonBID_MODIFIER_ALREADY_EXISTS CriterionErrorReason = "BID_MODIFIER_ALREADY_EXISTS"

	//
	// Feed Id is not allowed in these Location Groups.
	//
	CriterionErrorReasonFEED_ID_NOT_ALLOWED CriterionErrorReason = "FEED_ID_NOT_ALLOWED"

	//
	// The account may not use the requested criteria type. For example, some
	// accounts are restricted to keywords only.
	//
	CriterionErrorReasonACCOUNT_INELIGIBLE_FOR_CRITERIA_TYPE CriterionErrorReason = "ACCOUNT_INELIGIBLE_FOR_CRITERIA_TYPE"

	//
	// The requested criteria type cannot be used with campaign or ad group bidding strategy.
	//
	CriterionErrorReasonCRITERIA_TYPE_INVALID_FOR_BIDDING_STRATEGY CriterionErrorReason = "CRITERIA_TYPE_INVALID_FOR_BIDDING_STRATEGY"

	//
	// The Criterion is not allowed to be excluded.
	//
	CriterionErrorReasonCANNOT_EXCLUDE_CRITERION CriterionErrorReason = "CANNOT_EXCLUDE_CRITERION"

	//
	// The criterion is not allowed to be removed. For example, we cannot remove any
	// of the platform criterion.
	//
	CriterionErrorReasonCANNOT_REMOVE_CRITERION CriterionErrorReason = "CANNOT_REMOVE_CRITERION"

	//
	// The combined length of product dimension values of the product scope criterion is too long.
	//
	CriterionErrorReasonPRODUCT_SCOPE_TOO_LONG CriterionErrorReason = "PRODUCT_SCOPE_TOO_LONG"

	//
	// Product scope contains too many dimensions.
	//
	CriterionErrorReasonPRODUCT_SCOPE_TOO_MANY_DIMENSIONS CriterionErrorReason = "PRODUCT_SCOPE_TOO_MANY_DIMENSIONS"

	//
	// The combined length of product dimension values of the product partition criterion is too
	// long.
	//
	CriterionErrorReasonPRODUCT_PARTITION_TOO_LONG CriterionErrorReason = "PRODUCT_PARTITION_TOO_LONG"

	//
	// Product partition contains too many dimensions.
	//
	CriterionErrorReasonPRODUCT_PARTITION_TOO_MANY_DIMENSIONS CriterionErrorReason = "PRODUCT_PARTITION_TOO_MANY_DIMENSIONS"

	//
	// The product dimension is invalid (e.g. dimension contains illegal value, dimension type is
	// represented with wrong class, etc). Product dimension value can not contain "==" or "&+".
	//
	CriterionErrorReasonINVALID_PRODUCT_DIMENSION CriterionErrorReason = "INVALID_PRODUCT_DIMENSION"

	//
	// Product dimension type is either invalid for campaigns of this type or cannot be used in the
	// current context. BIDDING_CATEGORY_Lx and PRODUCT_TYPE_Lx product dimensions must be used in
	// ascending order of their levels: L1, L2, L3, L4, L5... The levels must be specified
	// sequentially and start from L1. Furthermore, an "others" product partition cannot be
	// subdivided with a dimension of the same type but of a higher level ("others"
	// BIDDING_CATEGORY_L3 can be subdivided with BRAND but not with BIDDING_CATEGORY_L4).
	//
	CriterionErrorReasonINVALID_PRODUCT_DIMENSION_TYPE CriterionErrorReason = "INVALID_PRODUCT_DIMENSION_TYPE"

	//
	// Bidding categories do not form a valid path in the Shopping bidding category taxonomy.
	//
	CriterionErrorReasonINVALID_PRODUCT_BIDDING_CATEGORY CriterionErrorReason = "INVALID_PRODUCT_BIDDING_CATEGORY"

	//
	// ShoppingSetting must be added to the campaign before ProductScope criteria can be added.
	//
	CriterionErrorReasonMISSING_SHOPPING_SETTING CriterionErrorReason = "MISSING_SHOPPING_SETTING"

	//
	// Matching function is invalid.
	//
	CriterionErrorReasonINVALID_MATCHING_FUNCTION CriterionErrorReason = "INVALID_MATCHING_FUNCTION"

	//
	// Filter parameters not allowed for location groups targeting.
	//
	CriterionErrorReasonLOCATION_FILTER_NOT_ALLOWED CriterionErrorReason = "LOCATION_FILTER_NOT_ALLOWED"

	//
	// Given location filter parameter is invalid for location groups targeting.
	//
	CriterionErrorReasonLOCATION_FILTER_INVALID CriterionErrorReason = "LOCATION_FILTER_INVALID"

	//
	// Criteria type cannot be associated with a campaign and its ad group(s) simultaneously.
	//
	CriterionErrorReasonCANNOT_ATTACH_CRITERIA_AT_CAMPAIGN_AND_ADGROUP CriterionErrorReason = "CANNOT_ATTACH_CRITERIA_AT_CAMPAIGN_AND_ADGROUP"

	CriterionErrorReasonUNKNOWN CriterionErrorReason = "UNKNOWN"
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
// Days of the week.
//
type DayOfWeek string

const (

	//
	// The day of week named Monday.
	//
	DayOfWeekMONDAY DayOfWeek = "MONDAY"

	//
	// The day of week named Tuesday.
	//
	DayOfWeekTUESDAY DayOfWeek = "TUESDAY"

	//
	// The day of week named Wednesday.
	//
	DayOfWeekWEDNESDAY DayOfWeek = "WEDNESDAY"

	//
	// The day of week named Thursday.
	//
	DayOfWeekTHURSDAY DayOfWeek = "THURSDAY"

	//
	// The day of week named Friday.
	//
	DayOfWeekFRIDAY DayOfWeek = "FRIDAY"

	//
	// The day of week named Saturday.
	//
	DayOfWeekSATURDAY DayOfWeek = "SATURDAY"

	//
	// The day of week named Sunday.
	//
	DayOfWeekSUNDAY DayOfWeek = "SUNDAY"
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
// Limits at various levels of the account.
//
type EntityCountLimitExceededReason string

const (

	//
	// Indicates that this request would exceed the number of allowed entities for the AdWords
	// account. The exact entity type and limit being checked can be inferred from
	// {@link #accountLimitType}.
	//
	EntityCountLimitExceededReasonACCOUNT_LIMIT EntityCountLimitExceededReason = "ACCOUNT_LIMIT"

	//
	// Indicates that this request would exceed the number of allowed entities in a Campaign.
	// The exact entity type and limit being checked can be inferred from
	// {@link #accountLimitType}, and the numeric id of the Campaign involved is given by
	// {@link #enclosingId}.
	//
	EntityCountLimitExceededReasonCAMPAIGN_LIMIT EntityCountLimitExceededReason = "CAMPAIGN_LIMIT"

	//
	// Indicates that this request would exceed the number of allowed entities in
	// an ad group.  The exact entity type and limit being checked can be
	// inferred from {@link #accountLimitType}, and the numeric id of the ad group
	// involved is given by {@link #enclosingId}.
	//
	EntityCountLimitExceededReasonADGROUP_LIMIT EntityCountLimitExceededReason = "ADGROUP_LIMIT"

	//
	// Indicates that this request would exceed the number of allowed entities in an ad group ad.
	// The exact entity type and limit being checked can be inferred from {@link #accountLimitType},
	// and the {@link #enclosingId} contains the ad group id followed by the ad id, separated by a
	// single comma (,).
	//
	EntityCountLimitExceededReasonAD_GROUP_AD_LIMIT EntityCountLimitExceededReason = "AD_GROUP_AD_LIMIT"

	//
	// Indicates that this request would exceed the number of allowed entities in an ad group
	// criterion.  The exact entity type and limit being checked can be inferred from
	// {@link #accountLimitType}, and the {@link #enclosingId} contains the ad group id followed by
	// the criterion id, separated by a single comma (,).
	//
	EntityCountLimitExceededReasonAD_GROUP_CRITERION_LIMIT EntityCountLimitExceededReason = "AD_GROUP_CRITERION_LIMIT"

	//
	// Indicates that this request would exceed the number of allowed entities in
	// this shared set.  The exact entity type and limit being checked can be
	// inferred from {@link #accountLimitType}, and the numeric id of the shared
	// set involved is given by {@link #enclosingId}.
	//
	EntityCountLimitExceededReasonSHARED_SET_LIMIT EntityCountLimitExceededReason = "SHARED_SET_LIMIT"

	//
	// Exceeds a limit related to a matching function.
	//
	EntityCountLimitExceededReasonMATCHING_FUNCTION_LIMIT EntityCountLimitExceededReason = "MATCHING_FUNCTION_LIMIT"

	//
	// Specific limit that has been exceeded is unknown (the client may be of an
	// older version than the server).
	//
	EntityCountLimitExceededReasonUNKNOWN EntityCountLimitExceededReason = "UNKNOWN"
)

type EntityNotFoundReason string

const (

	//
	// The specified id refered to an entity which either doesn't exist or is not accessible to the
	// customer. e.g. campaign belongs to another customer.
	//
	EntityNotFoundReasonINVALID_ID EntityNotFoundReason = "INVALID_ID"
)

//
// Operators that can be used in functions.
//
type FunctionOperator string

const (

	//
	// The IN operator.
	//
	FunctionOperatorIN FunctionOperator = "IN"

	//
	// The IDENTITY operator.
	//
	FunctionOperatorIDENTITY FunctionOperator = "IDENTITY"

	//
	// The EQUALS operator
	//
	FunctionOperatorEQUALS FunctionOperator = "EQUALS"

	//
	// Operator that takes two or more operands that are of type FunctionOperand
	// and checks that all the operands evaluate to true.
	// For functions related to ad formats, all the operands must be in lhsOperand.
	// Return ConstantOperand with Bool type.
	//
	FunctionOperatorAND FunctionOperator = "AND"

	//
	// Operator that returns true if the elements in lhsOperand contains any of the elements
	// in rhsOperands. Otherwise, return false.
	//
	FunctionOperatorCONTAINS_ANY FunctionOperator = "CONTAINS_ANY"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	FunctionOperatorUNKNOWN FunctionOperator = "UNKNOWN"
)

//
// The reasons for the target error.
//
type FunctionErrorReason string

const (

	//
	// The format of the function is not recognized as a supported function format.
	//
	FunctionErrorReasonINVALID_FUNCTION_FORMAT FunctionErrorReason = "INVALID_FUNCTION_FORMAT"

	//
	// Operand data types do not match.
	//
	FunctionErrorReasonDATA_TYPE_MISMATCH FunctionErrorReason = "DATA_TYPE_MISMATCH"

	//
	// The operands cannot be used together in a conjunction.
	//
	FunctionErrorReasonINVALID_CONJUNCTION_OPERANDS FunctionErrorReason = "INVALID_CONJUNCTION_OPERANDS"

	//
	// Invalid numer of Operands.
	//
	FunctionErrorReasonINVALID_NUMBER_OF_OPERANDS FunctionErrorReason = "INVALID_NUMBER_OF_OPERANDS"

	//
	// Operand Type not supported.
	//
	FunctionErrorReasonINVALID_OPERAND_TYPE FunctionErrorReason = "INVALID_OPERAND_TYPE"

	//
	// Operator not supported.
	//
	FunctionErrorReasonINVALID_OPERATOR FunctionErrorReason = "INVALID_OPERATOR"

	//
	// Request context type not supported.
	//
	FunctionErrorReasonINVALID_REQUEST_CONTEXT_TYPE FunctionErrorReason = "INVALID_REQUEST_CONTEXT_TYPE"

	//
	// The matching function is not allowed for call placeholders
	//
	FunctionErrorReasonINVALID_FUNCTION_FOR_CALL_PLACEHOLDER FunctionErrorReason = "INVALID_FUNCTION_FOR_CALL_PLACEHOLDER"

	//
	// The matching function is not allowed for the specified placeholder
	//
	FunctionErrorReasonINVALID_FUNCTION_FOR_PLACEHOLDER FunctionErrorReason = "INVALID_FUNCTION_FOR_PLACEHOLDER"

	//
	// Invalid operand.
	//
	FunctionErrorReasonINVALID_OPERAND FunctionErrorReason = "INVALID_OPERAND"

	//
	// Missing value for the constant operand.
	//
	FunctionErrorReasonMISSING_CONSTANT_OPERAND_VALUE FunctionErrorReason = "MISSING_CONSTANT_OPERAND_VALUE"

	//
	// The value of the constant operand is invalid.
	//
	FunctionErrorReasonINVALID_CONSTANT_OPERAND_VALUE FunctionErrorReason = "INVALID_CONSTANT_OPERAND_VALUE"

	//
	// Invalid function nesting.
	//
	FunctionErrorReasonINVALID_NESTING FunctionErrorReason = "INVALID_NESTING"

	//
	// The Feed ID was different from another Feed ID in the same function.
	//
	FunctionErrorReasonMULTIPLE_FEED_IDS_NOT_SUPPORTED FunctionErrorReason = "MULTIPLE_FEED_IDS_NOT_SUPPORTED"

	//
	// The matching function is invalid for use with a feed with a fixed schema.
	//
	FunctionErrorReasonINVALID_FUNCTION_FOR_FEED_WITH_FIXED_SCHEMA FunctionErrorReason = "INVALID_FUNCTION_FOR_FEED_WITH_FIXED_SCHEMA"

	//
	// Invalid attribute name.
	//
	FunctionErrorReasonINVALID_ATTRIBUTE_NAME FunctionErrorReason = "INVALID_ATTRIBUTE_NAME"

	FunctionErrorReasonUNKNOWN FunctionErrorReason = "UNKNOWN"
)

type GenderGenderType string

const (
	GenderGenderTypeGENDER_MALE GenderGenderType = "GENDER_MALE"

	GenderGenderTypeGENDER_FEMALE GenderGenderType = "GENDER_FEMALE"

	GenderGenderTypeGENDER_UNDETERMINED GenderGenderType = "GENDER_UNDETERMINED"
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
// Income percentile ranges.
//
type IncomeRangeIncomeRangeType string

const (

	//
	// Users for which income range is undetermined.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_UNDETERMINED IncomeRangeIncomeRangeType = "INCOME_RANGE_UNDETERMINED"

	//
	// Users in [0%, 50%) income percentile range.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_0_50 IncomeRangeIncomeRangeType = "INCOME_RANGE_0_50"

	//
	// Users in [50%, 60%) income percentile range.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_50_60 IncomeRangeIncomeRangeType = "INCOME_RANGE_50_60"

	//
	// Users in [60%, 70%) income percentile range.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_60_70 IncomeRangeIncomeRangeType = "INCOME_RANGE_60_70"

	//
	// Users in [70%, 80%) income percentile range.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_70_80 IncomeRangeIncomeRangeType = "INCOME_RANGE_70_80"

	//
	// Users in [80%, 90%) income percentile range.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_80_90 IncomeRangeIncomeRangeType = "INCOME_RANGE_80_90"

	//
	// Users in [90%, 100%] income percentile range.
	//
	IncomeRangeIncomeRangeTypeINCOME_RANGE_90_UP IncomeRangeIncomeRangeType = "INCOME_RANGE_90_UP"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	IncomeRangeIncomeRangeTypeUNKNOWN IncomeRangeIncomeRangeType = "UNKNOWN"
)

//
// Income tiers that specify the income bracket a household falls under. TIER_1
// belongs to the highest income bracket. The income bracket range associated with
// each tier is defined per country and computed based on income percentiles.
//
type IncomeTier string

const (
	IncomeTierUNKNOWN IncomeTier = "UNKNOWN"

	IncomeTierTIER_1 IncomeTier = "TIER_1"

	IncomeTierTIER_2 IncomeTier = "TIER_2"

	IncomeTierTIER_3 IncomeTier = "TIER_3"

	IncomeTierTIER_4 IncomeTier = "TIER_4"

	IncomeTierTIER_5 IncomeTier = "TIER_5"

	//
	// Bucket consisting of the bottom 5 tiers, specifying the bottom 50% of household
	// income zip codes.
	//
	IncomeTierTIER_6_TO_10 IncomeTier = "TIER_6_TO_10"
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
// Match type of a keyword. i.e. the way we match a keyword string with
// search queries.
//
type KeywordMatchType string

const (

	//
	// Exact match
	//
	KeywordMatchTypeEXACT KeywordMatchType = "EXACT"

	//
	// Phrase match
	//
	KeywordMatchTypePHRASE KeywordMatchType = "PHRASE"

	//
	// Broad match
	//
	KeywordMatchTypeBROAD KeywordMatchType = "BROAD"
)

//
// Enum that represents the different Targeting Status values for a Location criterion.
//
type LocationTargetingStatus string

const (

	//
	// The location is active.
	//
	LocationTargetingStatusACTIVE LocationTargetingStatus = "ACTIVE"

	//
	// The location is not available for targeting.
	//
	LocationTargetingStatusOBSOLETE LocationTargetingStatus = "OBSOLETE"

	//
	// The location is phasing out, it will marked obsolete soon.
	//
	LocationTargetingStatusPHASING_OUT LocationTargetingStatus = "PHASING_OUT"
)

//
// Minutes in an hour.  Currently only 0, 15, 30, and 45 are supported
//
type MinuteOfHour string

const (

	//
	// Zero minutes past hour.
	//
	MinuteOfHourZERO MinuteOfHour = "ZERO"

	//
	// Fifteen minutes past hour.
	//
	MinuteOfHourFIFTEEN MinuteOfHour = "FIFTEEN"

	//
	// Thirty minutes past hour.
	//
	MinuteOfHourTHIRTY MinuteOfHour = "THIRTY"

	//
	// Forty-five minutes past hour.
	//
	MinuteOfHourFORTY_FIVE MinuteOfHour = "FORTY_FIVE"
)

type MobileDeviceDeviceType string

const (
	MobileDeviceDeviceTypeDEVICE_TYPE_MOBILE MobileDeviceDeviceType = "DEVICE_TYPE_MOBILE"

	MobileDeviceDeviceTypeDEVICE_TYPE_TABLET MobileDeviceDeviceType = "DEVICE_TYPE_TABLET"
)

type NewEntityCreationErrorReason string

const (

	//
	// Do not set the id field while creating new entities.
	//
	NewEntityCreationErrorReasonCANNOT_SET_ID_FOR_ADD NewEntityCreationErrorReason = "CANNOT_SET_ID_FOR_ADD"

	//
	// Creating more than one entity with the same temp ID is not allowed.
	//
	NewEntityCreationErrorReasonDUPLICATE_TEMP_IDS NewEntityCreationErrorReason = "DUPLICATE_TEMP_IDS"

	//
	// Parent object with specified temp id failed validation, so no deep
	// validation will be done for this child entity.
	//
	NewEntityCreationErrorReasonTEMP_ID_ENTITY_HAD_ERRORS NewEntityCreationErrorReason = "TEMP_ID_ENTITY_HAD_ERRORS"
)

//
// The reasons for the validation error.
//
type NotEmptyErrorReason string

const (
	NotEmptyErrorReasonEMPTY_LIST NotEmptyErrorReason = "EMPTY_LIST"
)

//
// The reasons for the validation error.
//
type NullErrorReason string

const (

	//
	// Specified list/container must not contain any null elements
	//
	NullErrorReasonNULL_CONTENT NullErrorReason = "NULL_CONTENT"
)

//
// The operator type.
//
type OperatingSystemVersionOperatorType string

const (
	OperatingSystemVersionOperatorTypeGREATER_THAN_EQUAL_TO OperatingSystemVersionOperatorType = "GREATER_THAN_EQUAL_TO"

	OperatingSystemVersionOperatorTypeEQUAL_TO OperatingSystemVersionOperatorType = "EQUAL_TO"

	OperatingSystemVersionOperatorTypeUNKNOWN OperatingSystemVersionOperatorType = "UNKNOWN"
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
// The reasons for errors when using pagination.
//
type PagingErrorReason string

const (

	//
	// The start index value cannot be a negative number.
	//
	PagingErrorReasonSTART_INDEX_CANNOT_BE_NEGATIVE PagingErrorReason = "START_INDEX_CANNOT_BE_NEGATIVE"

	//
	// The number of results cannot be a negative number.
	//
	PagingErrorReasonNUMBER_OF_RESULTS_CANNOT_BE_NEGATIVE PagingErrorReason = "NUMBER_OF_RESULTS_CANNOT_BE_NEGATIVE"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	PagingErrorReasonUNKNOWN PagingErrorReason = "UNKNOWN"
)

//
// The possible types of parents.
//
type ParentParentType string

const (
	ParentParentTypePARENT_PARENT ParentParentType = "PARENT_PARENT"

	ParentParentTypePARENT_NOT_A_PARENT ParentParentType = "PARENT_NOT_A_PARENT"

	ParentParentTypePARENT_UNDETERMINED ParentParentType = "PARENT_UNDETERMINED"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ParentParentTypeUNKNOWN ParentParentType = "UNKNOWN"
)

//
// Categories to identify places of interest.
//
type PlacesOfInterestOperandCategory string

const (
	PlacesOfInterestOperandCategoryAIRPORT PlacesOfInterestOperandCategory = "AIRPORT"

	PlacesOfInterestOperandCategoryDOWNTOWN PlacesOfInterestOperandCategory = "DOWNTOWN"

	PlacesOfInterestOperandCategoryUNIVERSITY PlacesOfInterestOperandCategory = "UNIVERSITY"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	PlacesOfInterestOperandCategoryUNKNOWN PlacesOfInterestOperandCategory = "UNKNOWN"
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
// A canonical product condition.
//
type ProductCanonicalConditionCondition string

const (
	ProductCanonicalConditionConditionNEW ProductCanonicalConditionCondition = "NEW"

	ProductCanonicalConditionConditionUSED ProductCanonicalConditionCondition = "USED"

	ProductCanonicalConditionConditionREFURBISHED ProductCanonicalConditionCondition = "REFURBISHED"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ProductCanonicalConditionConditionUNKNOWN ProductCanonicalConditionCondition = "UNKNOWN"
)

//
// Type of product dimension.
//
type ProductDimensionType string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ProductDimensionTypeUNKNOWN ProductDimensionType = "UNKNOWN"

	ProductDimensionTypeBIDDING_CATEGORY_L1 ProductDimensionType = "BIDDING_CATEGORY_L1"

	ProductDimensionTypeBIDDING_CATEGORY_L2 ProductDimensionType = "BIDDING_CATEGORY_L2"

	ProductDimensionTypeBIDDING_CATEGORY_L3 ProductDimensionType = "BIDDING_CATEGORY_L3"

	ProductDimensionTypeBIDDING_CATEGORY_L4 ProductDimensionType = "BIDDING_CATEGORY_L4"

	ProductDimensionTypeBIDDING_CATEGORY_L5 ProductDimensionType = "BIDDING_CATEGORY_L5"

	ProductDimensionTypeBRAND ProductDimensionType = "BRAND"

	ProductDimensionTypeCANONICAL_CONDITION ProductDimensionType = "CANONICAL_CONDITION"

	ProductDimensionTypeCUSTOM_ATTRIBUTE_0 ProductDimensionType = "CUSTOM_ATTRIBUTE_0"

	ProductDimensionTypeCUSTOM_ATTRIBUTE_1 ProductDimensionType = "CUSTOM_ATTRIBUTE_1"

	ProductDimensionTypeCUSTOM_ATTRIBUTE_2 ProductDimensionType = "CUSTOM_ATTRIBUTE_2"

	ProductDimensionTypeCUSTOM_ATTRIBUTE_3 ProductDimensionType = "CUSTOM_ATTRIBUTE_3"

	ProductDimensionTypeCUSTOM_ATTRIBUTE_4 ProductDimensionType = "CUSTOM_ATTRIBUTE_4"

	ProductDimensionTypeOFFER_ID ProductDimensionType = "OFFER_ID"

	ProductDimensionTypePRODUCT_TYPE_L1 ProductDimensionType = "PRODUCT_TYPE_L1"

	ProductDimensionTypePRODUCT_TYPE_L2 ProductDimensionType = "PRODUCT_TYPE_L2"

	ProductDimensionTypePRODUCT_TYPE_L3 ProductDimensionType = "PRODUCT_TYPE_L3"

	ProductDimensionTypePRODUCT_TYPE_L4 ProductDimensionType = "PRODUCT_TYPE_L4"

	ProductDimensionTypePRODUCT_TYPE_L5 ProductDimensionType = "PRODUCT_TYPE_L5"

	ProductDimensionTypeCHANNEL ProductDimensionType = "CHANNEL"

	ProductDimensionTypeCHANNEL_EXCLUSIVITY ProductDimensionType = "CHANNEL_EXCLUSIVITY"
)

//
// The radius distance is expressed in either kilometers or miles.
//
type ProximityDistanceUnits string

const (

	//
	// The unit of distance is kilometer.
	//
	ProximityDistanceUnitsKILOMETERS ProximityDistanceUnits = "KILOMETERS"

	//
	// The unit of distance is mile.
	//
	ProximityDistanceUnitsMILES ProximityDistanceUnits = "MILES"
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
// The reasons for the validation error.
//
type RegionCodeErrorReason string

const (
	RegionCodeErrorReasonINVALID_REGION_CODE RegionCodeErrorReason = "INVALID_REGION_CODE"
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
// Channel specifies where the item is sold: online or in local stores.
//
type ShoppingProductChannel string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ShoppingProductChannelUNKNOWN ShoppingProductChannel = "UNKNOWN"

	//
	// The item is sold online.
	//
	ShoppingProductChannelONLINE ShoppingProductChannel = "ONLINE"

	//
	// The item is sold in local stores.
	//
	ShoppingProductChannelLOCAL ShoppingProductChannel = "LOCAL"
)

//
// Channel exclusivity specifies whether an item is sold exclusively in one channel
// or through multiple channels.
//
type ShoppingProductChannelExclusivity string

const (

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	ShoppingProductChannelExclusivityUNKNOWN ShoppingProductChannelExclusivity = "UNKNOWN"

	//
	// The item is sold through one channel only, either local stores or online as
	// indicated by its ShoppingProductChannel.
	//
	ShoppingProductChannelExclusivitySINGLE_CHANNEL ShoppingProductChannelExclusivity = "SINGLE_CHANNEL"

	//
	// The item is matched to its online or local stores counterpart, indicating it is
	// available for purchase in both ShoppingProductChannels.
	//
	ShoppingProductChannelExclusivityMULTI_CHANNEL ShoppingProductChannelExclusivity = "MULTI_CHANNEL"
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

//
// Membership status of the user list.
//
type CriterionUserListMembershipStatus string

const (

	//
	// Open status - list is accruing members and can be targeted to.
	//
	CriterionUserListMembershipStatusOPEN CriterionUserListMembershipStatus = "OPEN"

	//
	// Closed status - No new members being added. Can not be used for targeting a new campaign.
	// Existing campaigns can still work as long as the list is not removed as a targeting criteria.
	//
	CriterionUserListMembershipStatusCLOSED CriterionUserListMembershipStatus = "CLOSED"
)

//
// Operand value of {@link WebpageCondition}.
//
type WebpageConditionOperand string

const (

	//
	// Operand denoting a webpage URL targeting condition.
	// The operator {@link StringConditionOperator#CONTAINS} will be used for
	// such conditions.
	//
	WebpageConditionOperandURL WebpageConditionOperand = "URL"

	//
	// Operand denoting a webpage category targeting condition.
	// The operator {@link StringConditionOperator#EQUALS} will be used for
	// such conditions.
	//
	WebpageConditionOperandCATEGORY WebpageConditionOperand = "CATEGORY"

	//
	// Operand denoting a webpage title targeting condition.
	// The operator {@link StringConditionOperator#CONTAINS} will be used for
	// such conditions.
	//
	WebpageConditionOperandPAGE_TITLE WebpageConditionOperand = "PAGE_TITLE"

	//
	// Operand denoting a webpage content targeting condition.
	// The operator {@link StringConditionOperator#CONTAINS} will be used for
	// such conditions.
	//
	WebpageConditionOperandPAGE_CONTENT WebpageConditionOperand = "PAGE_CONTENT"

	//
	// Operand denoting a webpage custom label targeting condition.<br>
	// The operator {@link StringConditionOperator#EQUALS} will be used for such conditions.
	//
	WebpageConditionOperandCUSTOM_LABEL WebpageConditionOperand = "CUSTOM_LABEL"

	//
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	WebpageConditionOperandUNKNOWN WebpageConditionOperand = "UNKNOWN"
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

	Rval *CampaignCriterionPage `xml:"rval,omitempty"`
}

type Mutate struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutate"`

	//
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint DistinctIds">Elements in this field must have distinct IDs for following {@link Operator}s : SET, REMOVE.</span>
	// <span class="constraint NotEmpty">This field must contain at least one element.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint SupportedOperators">The following {@link Operator}s are supported: ADD, REMOVE, SET.</span>
	//
	Operations []*CampaignCriterionOperation `xml:"operations,omitempty"`
}

type MutateResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 mutateResponse"`

	Rval *CampaignCriterionReturnValue `xml:"rval,omitempty"`
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

	Rval *CampaignCriterionPage `xml:"rval,omitempty"`
}

type AdSchedule struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdSchedule"`

	*Criterion

	//
	// Day of the week the schedule applies to.
	// <span class="constraint Selectable">This field can be selected using the value "DayOfWeek".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE and SET.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	DayOfWeek *DayOfWeek `xml:"dayOfWeek,omitempty"`

	//
	// Starting hour in 24 hour time.
	// <span class="constraint Selectable">This field can be selected using the value "StartHour".</span>
	// <span class="constraint InRange">This field must be between 0 and 23, inclusive.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE and SET.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	StartHour int32 `xml:"startHour,omitempty"`

	//
	// Interval starts these minutes after the starting hour.
	// The value can be 0, 15, 30, and 45.
	// <span class="constraint Selectable">This field can be selected using the value "StartMinute".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE and SET.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	StartMinute *MinuteOfHour `xml:"startMinute,omitempty"`

	//
	// Ending hour in 24 hour time; <code>24</code> signifies end of the day.
	// <span class="constraint Selectable">This field can be selected using the value "EndHour".</span>
	// <span class="constraint InRange">This field must be between 0 and 24, inclusive.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE and SET.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	EndHour int32 `xml:"endHour,omitempty"`

	//
	// Interval ends these minutes after the ending hour.
	// The value can be 0, 15, 30, and 45.
	// <span class="constraint Selectable">This field can be selected using the value "EndMinute".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE and SET.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	EndMinute *MinuteOfHour `xml:"endMinute,omitempty"`
}

type Address struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Address"`

	//
	// Street address line 1; <code>null</code> if unknown.
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	StreetAddress string `xml:"streetAddress,omitempty"`

	//
	// Street address line 2; <code>null</code> if unknown.
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	StreetAddress2 string `xml:"streetAddress2,omitempty"`

	//
	// Name of the city; <code>null</code> if unknown.
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	CityName string `xml:"cityName,omitempty"`

	//
	// Province or state code; <code>null</code> if unknown.
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	ProvinceCode string `xml:"provinceCode,omitempty"`

	//
	// Province or state name; <code>null</code> if unknown.
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	ProvinceName string `xml:"provinceName,omitempty"`

	//
	// Postal code; <code>null</code> if unknown.
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	PostalCode string `xml:"postalCode,omitempty"`

	//
	// Country code; <code>null</code> if unknown.
	//
	CountryCode string `xml:"countryCode,omitempty"`
}

type AdxError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdxError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AdxErrorReason `xml:"reason,omitempty"`
}

type AgeRange struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AgeRange"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "AgeRangeType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	AgeRangeType *AgeRangeAgeRangeType `xml:"ageRangeType,omitempty"`
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

type CampaignCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignCriterion"`

	//
	// The campaign that the criterion is in.
	// <span class="constraint Selectable">This field can be selected using the value "CampaignId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	CampaignId int64 `xml:"campaignId,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "IsNegative".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	IsNegative bool `xml:"isNegative,omitempty"`

	//
	// The criterion part of the campaign criterion.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Criterion *Criterion `xml:"criterion,omitempty"`

	//
	// The modifier for bids when the criterion matches.
	//
	// <p>Valid modifier values range from {@code 0.1} to {@code 10.0}, with {@code 0.0} reserved
	// for opting out of platform criterion.
	// <p>To clear an existing bid modifier, specify {@code -1.0} (invalid for initial {@code ADD}
	// operations).
	// <span class="constraint Selectable">This field can be selected using the value "BidModifier".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint CampaignType">This field may not be set for campaign channel type SHOPPING with campaign channel subtype SHOPPING_UNIVERSAL_ADS.</span>
	//
	BidModifier float64 `xml:"bidModifier,omitempty"`

	//
	// The status for criteria.
	// <span class="constraint Selectable">This field can be selected using the value "CampaignCriterionStatus".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	CampaignCriterionStatus *CampaignCriterionCampaignCriterionStatus `xml:"campaignCriterionStatus,omitempty"`

	//
	// ID of the base campaign from which this draft/trial campaign criterion was created.
	// This field is only returned on get requests.
	// <span class="constraint Selectable">This field can be selected using the value "BaseCampaignId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	BaseCampaignId int64 `xml:"baseCampaignId,omitempty"`

	//
	// This Map provides a place to put new features and settings in older versions
	// of the AdWords API in the rare instance we need to introduce a new feature in
	// an older version.
	//
	// It is presently unused.  Do not set a value.
	//
	ForwardCompatibilityMap []*String_StringMapEntry `xml:"forwardCompatibilityMap,omitempty"`

	//
	// Indicates that this instance is a subtype of CampaignCriterion.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	CampaignCriterionType string `xml:"CampaignCriterion.Type,omitempty"`
}

type CampaignCriterionError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignCriterionError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *CampaignCriterionErrorReason `xml:"reason,omitempty"`
}

type CampaignCriterionOperation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignCriterionOperation"`

	*Operation

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *CampaignCriterion `xml:"operand,omitempty"`
}

type CampaignCriterionPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignCriterionPage"`

	*Page

	//
	// The result entries in this page.
	//
	Entries []*CampaignCriterion `xml:"entries,omitempty"`
}

type CampaignCriterionReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignCriterionReturnValue"`

	*ListReturnValue

	Value []*CampaignCriterion `xml:"value,omitempty"`

	//
	// List of partial failure errors.
	//
	PartialFailureErrors []*ApiError `xml:"partialFailureErrors,omitempty"`
}

type Carrier struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Carrier"`

	*Criterion

	//
	// Name of the carrier.
	// <span class="constraint Selectable">This field can be selected using the value "CarrierName".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Name string `xml:"name,omitempty"`

	//
	// Country code of the carrier.
	// Can be {@code null} if not applicable, e.g., for Carrier "Wifi".
	// <span class="constraint Selectable">This field can be selected using the value "CarrierCountryCode".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CountryCode string `xml:"countryCode,omitempty"`
}

type ClientTermsError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ClientTermsError"`

	*ApiError

	Reason *ClientTermsErrorReason `xml:"reason,omitempty"`
}

type ConstantOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ConstantOperand"`

	*FunctionArgumentOperand

	//
	// Type of constant in this operand.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Type_ *ConstantOperandConstantType `xml:"type,omitempty"`

	//
	// Units of constant in this operand.
	//
	Unit *ConstantOperandUnit `xml:"unit,omitempty"`

	//
	// Long value of the operand if it is a long type.
	//
	LongValue int64 `xml:"longValue,omitempty"`

	//
	// Boolean value of the operand if it is a boolean type.
	//
	BooleanValue bool `xml:"booleanValue,omitempty"`

	//
	// Double value of the operand if it is a double type.
	//
	DoubleValue float64 `xml:"doubleValue,omitempty"`

	//
	// String value of the operand if it is a string type.
	//
	StringValue string `xml:"stringValue,omitempty"`
}

type ContentLabel struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ContentLabel"`

	*Criterion

	//
	// Content label type
	// <span class="constraint Selectable">This field can be selected using the value "ContentLabelType".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	ContentLabelType *ContentLabelType `xml:"contentLabelType,omitempty"`
}

type Criterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Criterion"`

	//
	// ID of this criterion.
	// <span class="constraint Selectable">This field can be selected using the value "Id".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : SET, REMOVE.</span>
	//
	Id int64 `xml:"id,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "CriteriaType".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Type_ *CriterionType `xml:"type,omitempty"`

	//
	// Indicates that this instance is a subtype of Criterion.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	CriterionType string `xml:"Criterion.Type,omitempty"`
}

type CriterionError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionError"`

	*ApiError

	Reason *CriterionErrorReason `xml:"reason,omitempty"`
}

type CriterionParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionParameter"`

	//
	// Indicates that this instance is a subtype of CriterionParameter.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	CriterionParameterType string `xml:"CriterionParameter.Type,omitempty"`
}

type DatabaseError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DatabaseError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DatabaseErrorReason `xml:"reason,omitempty"`
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

type EntityCountLimitExceeded struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 EntityCountLimitExceeded"`

	*ApiError

	//
	// Specifies which level's limit was exceeded.
	//
	Reason *EntityCountLimitExceededReason `xml:"reason,omitempty"`

	//
	// Id of the entity whose limit was exceeded.
	//
	EnclosingId string `xml:"enclosingId,omitempty"`

	//
	// The limit which was exceeded.
	//
	Limit int32 `xml:"limit,omitempty"`

	//
	// The account limit type which was exceeded.
	//
	AccountLimitType string `xml:"accountLimitType,omitempty"`

	//
	// The count of existing entities.
	//
	ExistingCount int32 `xml:"existingCount,omitempty"`
}

type EntityNotFound struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 EntityNotFound"`

	*ApiError

	//
	// Reason for this error.
	//
	Reason *EntityNotFoundReason `xml:"reason,omitempty"`
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

type Function struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Function"`

	//
	// Operator for a function.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operator *FunctionOperator `xml:"operator,omitempty"`

	//
	// Operand on the LHS in the equation. This is also the operand to be used for
	// single operand expressions such as NOT.
	// <span class="constraint CollectionSize">The minimum size of this collection is 1.</span>
	//
	LhsOperand []*FunctionArgumentOperand `xml:"lhsOperand,omitempty"`

	//
	// Operand on the RHS of the equation.
	//
	RhsOperand []*FunctionArgumentOperand `xml:"rhsOperand,omitempty"`

	//
	// String representation of the {@code Function}.
	//
	// <p>For mutate actions, this field can be set instead of the {@code operator},
	// {@code lhsOperand}, and {@code rhsOperand} fields. This field will be parsed and used to
	// populate the other fields.
	//
	// <p>When {@code Function} objects are returned from get or mutate calls, this field contains the
	// string representation of the {@code Function}. Note that because multiple strings may map to
	// the same {@code Function} (whitespace and single versus double quotation marks, for example),
	// the value returned may not be identical to the string sent in the request.
	//
	FunctionString string `xml:"functionString,omitempty"`
}

type FunctionError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 FunctionError"`

	*ApiError

	//
	// The error reason represented by an enum
	//
	Reason *FunctionErrorReason `xml:"reason,omitempty"`
}

type Gender struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Gender"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "GenderType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	GenderType *GenderGenderType `xml:"genderType,omitempty"`
}

type GeoPoint struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 GeoPoint"`

	//
	// Micro degrees for the latitude.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	LatitudeInMicroDegrees int32 `xml:"latitudeInMicroDegrees,omitempty"`

	//
	// Micro degrees for the longitude.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	LongitudeInMicroDegrees int32 `xml:"longitudeInMicroDegrees,omitempty"`
}

type GeoTargetOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 GeoTargetOperand"`

	*FunctionArgumentOperand

	//
	// CriterionId of locations deciding the geographical scope.
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	//
	Locations []int64 `xml:"locations,omitempty"`
}

type IdError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 IdError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *IdErrorReason `xml:"reason,omitempty"`
}

type IncomeOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 IncomeOperand"`

	*FunctionArgumentOperand

	//
	// Income tier specifying an income bracket that a household falls under. Tier 1 belongs to the
	// highest income bracket.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Tier *IncomeTier `xml:"tier,omitempty"`
}

type IncomeRange struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 IncomeRange"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "IncomeRangeType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	IncomeRangeType *IncomeRangeIncomeRangeType `xml:"incomeRangeType,omitempty"`
}

type InternalApiError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 InternalApiError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *InternalApiErrorReason `xml:"reason,omitempty"`
}

type IpBlock struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 IpBlock"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "IpAddress".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	IpAddress string `xml:"ipAddress,omitempty"`
}

type Keyword struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Keyword"`

	*Criterion

	//
	// Text of this keyword (at most 80 characters and ten words).
	// <span class="constraint Selectable">This field can be selected using the value "KeywordText".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint MatchesRegex">Keyword text must not contain NUL (code point 0x0) characters. This is checked by the regular expression '[^\x00]*'.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Text string `xml:"text,omitempty"`

	//
	// Match type of this keyword.
	// <span class="constraint Selectable">This field can be selected using the value "KeywordMatchType".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	MatchType *KeywordMatchType `xml:"matchType,omitempty"`
}

type Language struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Language"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "LanguageCode".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Code string `xml:"code,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "LanguageName".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Name string `xml:"name,omitempty"`
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

type Location struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Location"`

	*Criterion

	//
	// Name of the location criterion. <b> Note:</b> This field is filterable only in
	// LocationCriterionService. If used as a filter, a location name cannot be greater than 300
	// characters.
	// <span class="constraint Selectable">This field can be selected using the value "LocationName".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	LocationName string `xml:"locationName,omitempty"`

	//
	// Display type of the location criterion.
	// <span class="constraint Selectable">This field can be selected using the value "DisplayType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DisplayType string `xml:"displayType,omitempty"`

	//
	// The targeting status of the location criterion.
	// <span class="constraint Selectable">This field can be selected using the value "TargetingStatus".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	TargetingStatus *LocationTargetingStatus `xml:"targetingStatus,omitempty"`

	//
	// Ordered list of parents of the location criterion.
	// <span class="constraint Selectable">This field can be selected using the value "ParentLocations".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ParentLocations []*Location `xml:"parentLocations,omitempty"`
}

type LocationExtensionOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 LocationExtensionOperand"`

	*FunctionArgumentOperand

	//
	// Distance in units specifying the radius around targeted locations.
	// Only long and double are supported constant types.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Radius *ConstantOperand `xml:"radius,omitempty"`

	//
	// Used to filter locations present in the location feed by location criterion id.
	//
	LocationId int64 `xml:"locationId,omitempty"`
}

type MobileAppCategory struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MobileAppCategory"`

	*Criterion

	//
	// ID of this mobile app category. A complete list of the available mobile app categories is
	// available <a href="/adwords/api/docs/appendix/mobileappcategories">here</a>.
	// <span class="constraint Selectable">This field can be selected using the value "MobileAppCategoryId".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	MobileAppCategoryId int32 `xml:"mobileAppCategoryId,omitempty"`

	//
	// Name of this mobile app category.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DisplayName string `xml:"displayName,omitempty"`
}

type MobileApplication struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MobileApplication"`

	*Criterion

	//
	// A string that uniquely identifies a mobile application to AdWords API. The format of this
	// string is "<code>{platform}-{platform_native_id}</code>", where <code>platform</code> is "1"
	// for iOS apps and "2" for Android apps, and where <code>platform_native_id</code> is the mobile
	// application identifier native to the corresponding platform.
	// For iOS, this native identifier is the 9 digit string that appears at the end of an App Store
	// URL (e.g., "476943146" for "Flood-It! 2" whose App Store link is
	// http://itunes.apple.com/us/app/flood-it!-2/id476943146).
	// For Android, this native identifier is the application's package name (e.g.,
	// "com.labpixies.colordrips" for "Color Drips" given Google Play link
	// https://play.google.com/store/apps/details?id=com.labpixies.colordrips).
	// A well formed app id for AdWords API would thus be "1-476943146" for iOS and
	// "2-com.labpixies.colordrips" for Android.
	// <span class="constraint Selectable">This field can be selected using the value "AppId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	AppId string `xml:"appId,omitempty"`

	//
	// Title of this mobile application.
	// <span class="constraint Selectable">This field can be selected using the value "DisplayName".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DisplayName string `xml:"displayName,omitempty"`
}

type MobileDevice struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MobileDevice"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "DeviceName".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DeviceName string `xml:"deviceName,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "ManufacturerName".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ManufacturerName string `xml:"manufacturerName,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "DeviceType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DeviceType *MobileDeviceDeviceType `xml:"deviceType,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "OperatingSystemName".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	OperatingSystemName string `xml:"operatingSystemName,omitempty"`
}

type NegativeCampaignCriterion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NegativeCampaignCriterion"`

	*CampaignCriterion
}

type NewEntityCreationError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NewEntityCreationError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NewEntityCreationErrorReason `xml:"reason,omitempty"`
}

type NotEmptyError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NotEmptyError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NotEmptyErrorReason `xml:"reason,omitempty"`
}

type NullError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NullError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NullErrorReason `xml:"reason,omitempty"`
}

type FunctionArgumentOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 FunctionArgumentOperand"`

	//
	// Indicates that this instance is a subtype of FunctionArgumentOperand.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	FunctionArgumentOperandType string `xml:"FunctionArgumentOperand.Type,omitempty"`
}

type OperatingSystemVersion struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 OperatingSystemVersion"`

	*Criterion

	//
	// The name of the operating system.
	// <span class="constraint Selectable">This field can be selected using the value "OperatingSystemName".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Name string `xml:"name,omitempty"`

	//
	// The OS Major Version number.
	// <span class="constraint Selectable">This field can be selected using the value "OsMajorVersion".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	OsMajorVersion int32 `xml:"osMajorVersion,omitempty"`

	//
	// The OS Minor Version number.
	// <span class="constraint Selectable">This field can be selected using the value "OsMinorVersion".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	OsMinorVersion int32 `xml:"osMinorVersion,omitempty"`

	//
	// The operator type.
	// <span class="constraint Selectable">This field can be selected using the value "OperatorType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	OperatorType *OperatingSystemVersionOperatorType `xml:"operatorType,omitempty"`
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

type PagingError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 PagingError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *PagingErrorReason `xml:"reason,omitempty"`
}

type Parent struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Parent"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "ParentType".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ParentType *ParentParentType `xml:"parentType,omitempty"`
}

type Placement struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Placement"`

	*Criterion

	//
	// Url of the placement.
	//
	// <p>For example, "http://www.domain.com".
	// <span class="constraint Selectable">This field can be selected using the value "PlacementUrl".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Url string `xml:"url,omitempty"`
}

type PlacesOfInterestOperand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 PlacesOfInterestOperand"`

	*FunctionArgumentOperand

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Category *PlacesOfInterestOperandCategory `xml:"category,omitempty"`
}

type Platform struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Platform"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "PlatformName".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	PlatformName string `xml:"platformName,omitempty"`
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

type ProductAdwordsGrouping struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductAdwordsGrouping"`

	*ProductDimension

	//
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	Value string `xml:"value,omitempty"`
}

type ProductAdwordsLabels struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductAdwordsLabels"`

	*ProductDimension

	//
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	Value string `xml:"value,omitempty"`
}

type ProductBiddingCategory struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductBiddingCategory"`

	*ProductDimension

	//
	// Dimension type of the category. Indicates the level of the category in the taxonomy.
	// <span class="constraint Filterable">This field can be filtered on using the value "ParentDimensionType".</span>
	// <span class="constraint OneOf">The value must be one of {BIDDING_CATEGORY_L1, BIDDING_CATEGORY_L2, BIDDING_CATEGORY_L3, BIDDING_CATEGORY_L4, BIDDING_CATEGORY_L5}.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Type_ *ProductDimensionType `xml:"type,omitempty"`

	//
	// ID of the product category.
	// <span class="constraint Filterable">This field can be filtered on using the value "ParentDimensionId".</span>
	//
	Value int64 `xml:"value,omitempty"`
}

type ProductBrand struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductBrand"`

	*ProductDimension

	//
	// <span class="constraint StringLength">This string must not be empty, (trimmed).</span>
	//
	Value string `xml:"value,omitempty"`
}

type ProductCanonicalCondition struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductCanonicalCondition"`

	*ProductDimension

	Condition *ProductCanonicalConditionCondition `xml:"condition,omitempty"`
}

type ProductChannel struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductChannel"`

	*ProductDimension

	Channel *ShoppingProductChannel `xml:"channel,omitempty"`
}

type ProductChannelExclusivity struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductChannelExclusivity"`

	*ProductDimension

	ChannelExclusivity *ShoppingProductChannelExclusivity `xml:"channelExclusivity,omitempty"`
}

type ProductLegacyCondition struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductLegacyCondition"`

	*ProductDimension

	Value string `xml:"value,omitempty"`
}

type ProductCustomAttribute struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductCustomAttribute"`

	*ProductDimension

	//
	// Dimension type of the custom attribute. Indicates the index of the custom attribute.
	// <span class="constraint OneOf">The value must be one of {CUSTOM_ATTRIBUTE_0, CUSTOM_ATTRIBUTE_1, CUSTOM_ATTRIBUTE_2, CUSTOM_ATTRIBUTE_3, CUSTOM_ATTRIBUTE_4}.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Type_ *ProductDimensionType `xml:"type,omitempty"`

	//
	// <span class="constraint StringLength">This string must not be empty, (trimmed).</span>
	//
	Value string `xml:"value,omitempty"`
}

type ProductDimension struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductDimension"`

	//
	// Indicates that this instance is a subtype of ProductDimension.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	ProductDimensionType string `xml:"ProductDimension.Type,omitempty"`
}

type ProductOfferId struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductOfferId"`

	*ProductDimension

	//
	// <span class="constraint StringLength">This string must not be empty, (trimmed).</span>
	//
	Value string `xml:"value,omitempty"`
}

type ProductScope struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductScope"`

	*Criterion

	//
	// <span class="constraint Selectable">This field can be selected using the value "Dimensions".</span>
	// <span class="constraint NotEmptyForOperators">This field must contain at least one element when it is contained within {@link Operator}s: ADD.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Dimensions []*ProductDimension `xml:"dimensions,omitempty"`
}

type ProductType struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductType"`

	*ProductDimension

	//
	// Dimension type of the product type. Indicates the level of the product type.
	// <span class="constraint OneOf">The value must be one of {PRODUCT_TYPE_L1, PRODUCT_TYPE_L2, PRODUCT_TYPE_L3, PRODUCT_TYPE_L4, PRODUCT_TYPE_L5}.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Type_ *ProductDimensionType `xml:"type,omitempty"`

	//
	// <span class="constraint StringLength">This string must not be empty, (trimmed).</span>
	//
	Value string `xml:"value,omitempty"`
}

type ProductTypeFull struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ProductTypeFull"`

	*ProductDimension

	//
	// <span class="constraint StringLength">This string must not be empty.</span>
	//
	Value string `xml:"value,omitempty"`
}

type Proximity struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Proximity"`

	*Criterion

	//
	// Latitude and longitude.
	// <span class="constraint Selectable">This field can be selected using the value "GeoPoint".</span>
	//
	GeoPoint *GeoPoint `xml:"geoPoint,omitempty"`

	//
	// Radius distance units.
	// <span class="constraint Selectable">This field can be selected using the value "RadiusDistanceUnits".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	RadiusDistanceUnits *ProximityDistanceUnits `xml:"radiusDistanceUnits,omitempty"`

	//
	// Radius expressed in distance units.
	// <span class="constraint Selectable">This field can be selected using the value "RadiusInUnits".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	RadiusInUnits float64 `xml:"radiusInUnits,omitempty"`

	//
	// Full address; <code>null</code> if unknonwn.
	// <span class="constraint Selectable">This field can be selected using the value "Address".</span>
	//
	Address *Address `xml:"address,omitempty"`
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

type RegionCodeError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 RegionCodeError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *RegionCodeErrorReason `xml:"reason,omitempty"`
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

type LocationGroups struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 LocationGroups"`

	*Criterion

	//
	// Feed to be used for targeting around locations. This is required for distance targets.
	// <span class="constraint Selectable">This field can be selected using the value "FeedId".</span>
	//
	FeedId int64 `xml:"feedId,omitempty"`

	//
	// Matching function to filter out locations targeted by the criteria.
	//
	// This allows advertisers to target based on the semantics of the location.
	// <span class="constraint Selectable">This field can be selected using the value "MatchingFunction".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	MatchingFunction *Function `xml:"matchingFunction,omitempty"`
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

type String_StringMapEntry struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 String_StringMapEntry"`

	Key string `xml:"key,omitempty"`

	Value string `xml:"value,omitempty"`
}

type UnknownProductDimension struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 UnknownProductDimension"`

	*ProductDimension
}

type CriterionUserInterest struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionUserInterest"`

	*Criterion

	//
	// Id of this user interest. This is a required field.
	// <span class="constraint Selectable">This field can be selected using the value "UserInterestId".</span>
	//
	UserInterestId int64 `xml:"userInterestId,omitempty"`

	//
	// Parent Id of this user interest.
	// <span class="constraint Selectable">This field can be selected using the value "UserInterestParentId".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	UserInterestParentId int64 `xml:"userInterestParentId,omitempty"`

	//
	// Name of this user interest.
	// <span class="constraint Selectable">This field can be selected using the value "UserInterestName".</span>
	//
	UserInterestName string `xml:"userInterestName,omitempty"`
}

type CriterionUserList struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionUserList"`

	*Criterion

	//
	// Id of this user list. This is a required field.
	// <span class="constraint Selectable">This field can be selected using the value "UserListId".</span>
	//
	UserListId int64 `xml:"userListId,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "UserListName".</span>
	//
	UserListName string `xml:"userListName,omitempty"`

	//
	// <span class="constraint Selectable">This field can be selected using the value "UserListMembershipStatus".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	UserListMembershipStatus *CriterionUserListMembershipStatus `xml:"userListMembershipStatus,omitempty"`

	//
	// Determines whether a user list is eligible for targeting in the google.com
	// (search) network.
	// <span class="constraint Selectable">This field can be selected using the value "UserListEligibleForSearch".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	UserListEligibleForSearch bool `xml:"userListEligibleForSearch,omitempty"`

	//
	// Determines whether a user list is eligible for targeting in the display network.
	// <span class="constraint Selectable">This field can be selected using the value "UserListEligibleForDisplay".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	UserListEligibleForDisplay bool `xml:"userListEligibleForDisplay,omitempty"`
}

type Vertical struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Vertical"`

	*Criterion

	//
	// Id of this vertical.
	// <span class="constraint Selectable">This field can be selected using the value "VerticalId".</span>
	//
	VerticalId int64 `xml:"verticalId,omitempty"`

	//
	// Id of the parent of this vertical.
	// <span class="constraint Selectable">This field can be selected using the value "VerticalParentId".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	VerticalParentId int64 `xml:"verticalParentId,omitempty"`

	//
	// The category to target or exclude. Each subsequent element in the array
	// describes a more specific sub-category. For example,
	// <code>{"Pets &amp; Animals", "Pets", "Dogs"}</code> represents the "Pets &amp;
	// Animals/Pets/Dogs" category. A complete list of available vertical categories
	// is available <a href="/adwords/api/docs/appendix/verticals">here</a>
	// This field is required and must not be empty.
	// <span class="constraint Selectable">This field can be selected using the value "Path".</span>
	//
	Path []string `xml:"path,omitempty"`
}

type Webpage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Webpage"`

	*Criterion

	//
	// The webpage criterion parameter.
	// <span class="constraint Selectable">This field can be selected using the value "Parameter".</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	Parameter *WebpageParameter `xml:"parameter,omitempty"`

	//
	// Keywordless criteria coverage - Computed percentage of website coverage based on the
	// website target, negative website targets and negative keywords in the ad group and campaign.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CriteriaCoverage float64 `xml:"criteriaCoverage,omitempty"`

	//
	// Keywordless criteria samples - List of sample urls that matches with the website target.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	CriteriaSamples []string `xml:"criteriaSamples,omitempty"`
}

type WebpageCondition struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 WebpageCondition"`

	//
	// Operand of webpage targeting condition.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *WebpageConditionOperand `xml:"operand,omitempty"`

	//
	// Argument of the webpage targeting condition.
	// <span class="constraint MustNotContain">This string must not contain a substring that matches the regular expression '\*|\>\>|\=\=|\&\+'</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint StringLength">The length of this string should be between 1 and 2048, inclusive.</span>
	//
	Argument string `xml:"argument,omitempty"`
}

type WebpageParameter struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 WebpageParameter"`

	*CriterionParameter

	//
	// The name of the criterion that is defined by this parameter.
	//
	// <p>This name value will be used for identifying, sorting and filtering
	// criteria with this type of parameters. For criteria with simpler
	// parameters, such as keywords and placements, the parameter value itself
	// is used for identification, sorting and filtering.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint StringLength">The length of this string should be between 1 and 2048, inclusive.</span>
	//
	CriterionName string `xml:"criterionName,omitempty"`

	//
	// Conditions, or logical expressions, for webpage targeting.
	//
	// <p>The list of webpage targeting conditions are {@code and}-ed together
	// when evaluated for targeting. A {@code null} list of conditions means that
	// all webpages of the campaign's website are targeted.</p>
	// <span class="constraint CollectionSize">The maximum size of this collection is 3.</span>
	// <span class="constraint ContentsDistinct">This field must contain distinct elements.</span>
	//
	Conditions []*WebpageCondition `xml:"conditions,omitempty"`
}

type YouTubeChannel struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 YouTubeChannel"`

	*Criterion

	//
	// The YouTube uploader channel id or the channel code of a YouTube content channel.
	// <p>The uploader channel id can be obtained from the YouTube id-based URL. For example, in
	// <code>https://www.youtube.com/channel/UCEN58iXQg82TXgsDCjWqIkg</code> the channel id is
	// <code>UCEN58iXQg82TXgsDCjWqIkg</code>
	// <p>For more information see: https://support.google.com/youtube/answer/6180214
	// <span class="constraint Selectable">This field can be selected using the value "ChannelId".</span>
	//
	ChannelId string `xml:"channelId,omitempty"`

	//
	// The public name for a YouTube user channel.
	// <span class="constraint Selectable">This field can be selected using the value "ChannelName".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	ChannelName string `xml:"channelName,omitempty"`
}

type YouTubeVideo struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 YouTubeVideo"`

	*Criterion

	//
	// YouTube video id as it appears on the YouTube watch page.
	// <span class="constraint Selectable">This field can be selected using the value "VideoId".</span>
	//
	VideoId string `xml:"videoId,omitempty"`

	//
	// Name of the video.
	// <span class="constraint Selectable">This field can be selected using the value "VideoName".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	VideoName string `xml:"videoName,omitempty"`
}

type CampaignCriterionServiceInterface struct {
	client *SOAPClient
}

func NewCampaignCriterionServiceInterface(url string, tls bool, auth *BasicAuth) *CampaignCriterionServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &CampaignCriterionServiceInterface{
		client: client,
	}
}

func NewCampaignCriterionServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *CampaignCriterionServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &CampaignCriterionServiceInterface{
		client: client,
	}
}

func (service *CampaignCriterionServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *CampaignCriterionServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Gets campaign criteria.

   @param serviceSelector The selector specifying the {@link CampaignCriterion}s to return.
   @return A list of campaign criteria.
   @throws ApiException when there is at least one error with the request.
*/
func (service *CampaignCriterionServiceInterface) Get(request *Get) (*GetResponse, error) {
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
   Adds, removes or updates campaign criteria.

   @param operations The operations to apply.
   @return The added campaign criteria (without any optional parts).
   @throws ApiException when there is at least one error with the request.
*/
func (service *CampaignCriterionServiceInterface) Mutate(request *Mutate) (*MutateResponse, error) {
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
   Returns the list of campaign criteria that match the query.

   @param query The SQL-like AWQL query string.
   @return A list of campaign criteria.
   @throws ApiException if problems occur while parsing the query or fetching campaign criteria.
*/
func (service *CampaignCriterionServiceInterface) Query(request *Query) (*QueryResponse, error) {
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
