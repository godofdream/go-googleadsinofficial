package OfflineDataUploadService

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
// The reasons for the target error.
//
type AdErrorReason string

const (

	//
	// Ad customizers are not supported for ad type.
	//
	AdErrorReasonAD_CUSTOMIZERS_NOT_SUPPORTED_FOR_AD_TYPE AdErrorReason = "AD_CUSTOMIZERS_NOT_SUPPORTED_FOR_AD_TYPE"

	//
	// Estimating character sizes the string is too long.
	//
	AdErrorReasonAPPROXIMATELY_TOO_LONG AdErrorReason = "APPROXIMATELY_TOO_LONG"

	//
	// Estimating character sizes the string is too short.
	//
	AdErrorReasonAPPROXIMATELY_TOO_SHORT AdErrorReason = "APPROXIMATELY_TOO_SHORT"

	//
	// There is a problem with the snippet.
	//
	AdErrorReasonBAD_SNIPPET AdErrorReason = "BAD_SNIPPET"

	//
	// Cannot modify an ad.
	//
	AdErrorReasonCANNOT_MODIFY_AD AdErrorReason = "CANNOT_MODIFY_AD"

	//
	// business name and url cannot be set at the same time
	//
	AdErrorReasonCANNOT_SET_BUSINESS_NAME_IF_URL_SET AdErrorReason = "CANNOT_SET_BUSINESS_NAME_IF_URL_SET"

	//
	// The specified field is incompatible with this ad's type or settings.
	//
	AdErrorReasonCANNOT_SET_FIELD AdErrorReason = "CANNOT_SET_FIELD"

	//
	// Cannot set field when originAdId is set.
	//
	AdErrorReasonCANNOT_SET_FIELD_WITH_ORIGIN_AD_ID_SET AdErrorReason = "CANNOT_SET_FIELD_WITH_ORIGIN_AD_ID_SET"

	//
	// Cannot set field when an existing ad id is set for sharing.
	//
	AdErrorReasonCANNOT_SET_FIELD_WITH_AD_ID_SET_FOR_SHARING AdErrorReason = "CANNOT_SET_FIELD_WITH_AD_ID_SET_FOR_SHARING"

	//
	// Cannot set allowFlexibleColor false if no color is provided by user.
	//
	AdErrorReasonCANNOT_SET_ALLOW_FLEXIBLE_COLOR_FALSE AdErrorReason = "CANNOT_SET_ALLOW_FLEXIBLE_COLOR_FALSE"

	//
	// When user select native, no color control is allowed because we will always respect publisher
	// color for native format serving.
	//
	AdErrorReasonCANNOT_SET_COLOR_CONTROL_WHEN_NATIVE_FORMAT_SETTING AdErrorReason = "CANNOT_SET_COLOR_CONTROL_WHEN_NATIVE_FORMAT_SETTING"

	//
	// Cannot specify a url for the ad type
	//
	AdErrorReasonCANNOT_SET_URL AdErrorReason = "CANNOT_SET_URL"

	//
	// Cannot specify a tracking or mobile url without also setting final urls
	//
	AdErrorReasonCANNOT_SET_WITHOUT_FINAL_URLS AdErrorReason = "CANNOT_SET_WITHOUT_FINAL_URLS"

	//
	// Cannot specify a legacy url and a final url simultaneously
	//
	AdErrorReasonCANNOT_SET_WITH_FINAL_URLS AdErrorReason = "CANNOT_SET_WITH_FINAL_URLS"

	//
	// Cannot specify a legacy url and a tracking url template simultaneously in a DSA.
	//
	AdErrorReasonCANNOT_SET_WITH_TRACKING_URL_TEMPLATE AdErrorReason = "CANNOT_SET_WITH_TRACKING_URL_TEMPLATE"

	//
	// Cannot specify a urls in UrlData and in template fields simultaneously.
	//
	AdErrorReasonCANNOT_SET_WITH_URL_DATA AdErrorReason = "CANNOT_SET_WITH_URL_DATA"

	//
	// This operator cannot be used with a subclass of Ad.
	//
	AdErrorReasonCANNOT_USE_AD_SUBCLASS_FOR_OPERATOR AdErrorReason = "CANNOT_USE_AD_SUBCLASS_FOR_OPERATOR"

	//
	// Customer is not approved for mobile ads.
	//
	AdErrorReasonCUSTOMER_NOT_APPROVED_MOBILEADS AdErrorReason = "CUSTOMER_NOT_APPROVED_MOBILEADS"

	//
	// Customer is not approved for 3PAS richmedia ads.
	//
	AdErrorReasonCUSTOMER_NOT_APPROVED_THIRDPARTY_ADS AdErrorReason = "CUSTOMER_NOT_APPROVED_THIRDPARTY_ADS"

	//
	// Customer is not approved for 3PAS redirect richmedia (Ad Exchange) ads.
	//
	AdErrorReasonCUSTOMER_NOT_APPROVED_THIRDPARTY_REDIRECT_ADS AdErrorReason = "CUSTOMER_NOT_APPROVED_THIRDPARTY_REDIRECT_ADS"

	//
	// Not an eligible customer
	//
	AdErrorReasonCUSTOMER_NOT_ELIGIBLE AdErrorReason = "CUSTOMER_NOT_ELIGIBLE"

	//
	// Customer is not eligible for updating beacon url
	//
	AdErrorReasonCUSTOMER_NOT_ELIGIBLE_FOR_UPDATING_BEACON_URL AdErrorReason = "CUSTOMER_NOT_ELIGIBLE_FOR_UPDATING_BEACON_URL"

	//
	// There already exists an ad with the same dimensions in the union.
	//
	AdErrorReasonDIMENSION_ALREADY_IN_UNION AdErrorReason = "DIMENSION_ALREADY_IN_UNION"

	//
	// Ad's dimension must be set before setting union dimension.
	//
	AdErrorReasonDIMENSION_MUST_BE_SET AdErrorReason = "DIMENSION_MUST_BE_SET"

	//
	// Ad's dimension must be included in the union dimensions.
	//
	AdErrorReasonDIMENSION_NOT_IN_UNION AdErrorReason = "DIMENSION_NOT_IN_UNION"

	//
	// Display Url cannot be specified (applies to Ad Exchange Ads)
	//
	AdErrorReasonDISPLAY_URL_CANNOT_BE_SPECIFIED AdErrorReason = "DISPLAY_URL_CANNOT_BE_SPECIFIED"

	//
	// Telephone number contains invalid characters or invalid format.
	// Please re-enter your number using digits (0-9), dashes (-), and parentheses only.
	//
	AdErrorReasonDOMESTIC_PHONE_NUMBER_FORMAT AdErrorReason = "DOMESTIC_PHONE_NUMBER_FORMAT"

	//
	// Emergency telephone numbers are not allowed.
	// Please enter a valid domestic phone number to connect customers to your business.
	//
	AdErrorReasonEMERGENCY_PHONE_NUMBER AdErrorReason = "EMERGENCY_PHONE_NUMBER"

	//
	// A required field was not specified or is an empty string.
	//
	AdErrorReasonEMPTY_FIELD AdErrorReason = "EMPTY_FIELD"

	//
	// A feed attribute referenced in an ad customizer tag is not in the ad customizer mapping for
	// the feed.
	//
	AdErrorReasonFEED_ATTRIBUTE_MUST_HAVE_MAPPING_FOR_TYPE_ID AdErrorReason = "FEED_ATTRIBUTE_MUST_HAVE_MAPPING_FOR_TYPE_ID"

	//
	// The ad customizer field mapping for the feed attribute does not match the expected field
	// type.
	//
	AdErrorReasonFEED_ATTRIBUTE_MAPPING_TYPE_MISMATCH AdErrorReason = "FEED_ATTRIBUTE_MAPPING_TYPE_MISMATCH"

	//
	// The use of ad customizer tags in the ad text is disallowed. Details in trigger.
	//
	AdErrorReasonILLEGAL_AD_CUSTOMIZER_TAG_USE AdErrorReason = "ILLEGAL_AD_CUSTOMIZER_TAG_USE"

	//
	// Tags of the form {PH_x}, where x is a number, are disallowed in ad text.
	//
	AdErrorReasonILLEGAL_TAG_USE AdErrorReason = "ILLEGAL_TAG_USE"

	//
	// The dimensions of the ad are specified or derived in multiple ways and are not consistent.
	//
	AdErrorReasonINCONSISTENT_DIMENSIONS AdErrorReason = "INCONSISTENT_DIMENSIONS"

	//
	// The status cannot differ among template ads of the same union.
	//
	AdErrorReasonINCONSISTENT_STATUS_IN_TEMPLATE_UNION AdErrorReason = "INCONSISTENT_STATUS_IN_TEMPLATE_UNION"

	//
	// The length of the string is not valid.
	//
	AdErrorReasonINCORRECT_LENGTH AdErrorReason = "INCORRECT_LENGTH"

	//
	// The ad is ineligible for upgrade.
	//
	AdErrorReasonINELIGIBLE_FOR_UPGRADE AdErrorReason = "INELIGIBLE_FOR_UPGRADE"

	//
	// User cannot create mobile ad for countries targeted in specified campaign.
	//
	AdErrorReasonINVALID_AD_ADDRESS_CAMPAIGN_TARGET AdErrorReason = "INVALID_AD_ADDRESS_CAMPAIGN_TARGET"

	//
	// Invalid Ad type. A specific type of Ad is required.
	//
	AdErrorReasonINVALID_AD_TYPE AdErrorReason = "INVALID_AD_TYPE"

	//
	// Headline, description or phone cannot be present when creating mobile image ad.
	//
	AdErrorReasonINVALID_ATTRIBUTES_FOR_MOBILE_IMAGE AdErrorReason = "INVALID_ATTRIBUTES_FOR_MOBILE_IMAGE"

	//
	// Image cannot be present when creating mobile text ad.
	//
	AdErrorReasonINVALID_ATTRIBUTES_FOR_MOBILE_TEXT AdErrorReason = "INVALID_ATTRIBUTES_FOR_MOBILE_TEXT"

	//
	// Invalid call to action text.
	//
	AdErrorReasonINVALID_CALL_TO_ACTION_TEXT AdErrorReason = "INVALID_CALL_TO_ACTION_TEXT"

	//
	// Invalid character in URL.
	//
	AdErrorReasonINVALID_CHARACTER_FOR_URL AdErrorReason = "INVALID_CHARACTER_FOR_URL"

	//
	// Creative's country code is not valid.
	//
	AdErrorReasonINVALID_COUNTRY_CODE AdErrorReason = "INVALID_COUNTRY_CODE"

	//
	// Invalid use of Dynamic Search Ads tags ({lpurl} etc.)
	//
	AdErrorReasonINVALID_DSA_URL_TAG AdErrorReason = "INVALID_DSA_URL_TAG"

	//
	// Invalid use of Expanded Dynamic Search Ads tags ({lpurl} etc.)
	//
	AdErrorReasonINVALID_EXPANDED_DYNAMIC_SEARCH_AD_TAG AdErrorReason = "INVALID_EXPANDED_DYNAMIC_SEARCH_AD_TAG"

	//
	// An input error whose real reason was not properly mapped (should not happen).
	//
	AdErrorReasonINVALID_INPUT AdErrorReason = "INVALID_INPUT"

	//
	// An invalid markup language was entered.
	//
	AdErrorReasonINVALID_MARKUP_LANGUAGE AdErrorReason = "INVALID_MARKUP_LANGUAGE"

	//
	// An invalid mobile carrier was entered.
	//
	AdErrorReasonINVALID_MOBILE_CARRIER AdErrorReason = "INVALID_MOBILE_CARRIER"

	//
	// Specified mobile carriers target a country not targeted by the campaign.
	//
	AdErrorReasonINVALID_MOBILE_CARRIER_TARGET AdErrorReason = "INVALID_MOBILE_CARRIER_TARGET"

	//
	// Wrong number of elements for given element type
	//
	AdErrorReasonINVALID_NUMBER_OF_ELEMENTS AdErrorReason = "INVALID_NUMBER_OF_ELEMENTS"

	//
	// The format of the telephone number is incorrect.
	// Please re-enter the number using the correct format.
	//
	AdErrorReasonINVALID_PHONE_NUMBER_FORMAT AdErrorReason = "INVALID_PHONE_NUMBER_FORMAT"

	//
	// The certified vendor format id is incorrect.
	//
	AdErrorReasonINVALID_RICH_MEDIA_CERTIFIED_VENDOR_FORMAT_ID AdErrorReason = "INVALID_RICH_MEDIA_CERTIFIED_VENDOR_FORMAT_ID"

	//
	// The template ad data contains validation errors.
	//
	AdErrorReasonINVALID_TEMPLATE_DATA AdErrorReason = "INVALID_TEMPLATE_DATA"

	//
	// The template field doesn't have have the correct type.
	//
	AdErrorReasonINVALID_TEMPLATE_ELEMENT_FIELD_TYPE AdErrorReason = "INVALID_TEMPLATE_ELEMENT_FIELD_TYPE"

	//
	// Invalid template id.
	//
	AdErrorReasonINVALID_TEMPLATE_ID AdErrorReason = "INVALID_TEMPLATE_ID"

	//
	// After substituting replacement strings, the line is too wide.
	//
	AdErrorReasonLINE_TOO_WIDE AdErrorReason = "LINE_TOO_WIDE"

	//
	// The feed referenced must have ad customizer mapping to be used in a customizer tag.
	//
	AdErrorReasonMISSING_AD_CUSTOMIZER_MAPPING AdErrorReason = "MISSING_AD_CUSTOMIZER_MAPPING"

	//
	// Missing address component in template element address field.
	//
	AdErrorReasonMISSING_ADDRESS_COMPONENT AdErrorReason = "MISSING_ADDRESS_COMPONENT"

	//
	// An ad name must be entered.
	//
	AdErrorReasonMISSING_ADVERTISEMENT_NAME AdErrorReason = "MISSING_ADVERTISEMENT_NAME"

	//
	// Business name must be entered.
	//
	AdErrorReasonMISSING_BUSINESS_NAME AdErrorReason = "MISSING_BUSINESS_NAME"

	//
	// Description (line 2) must be entered.
	//
	AdErrorReasonMISSING_DESCRIPTION1 AdErrorReason = "MISSING_DESCRIPTION1"

	//
	// Description (line 3) must be entered.
	//
	AdErrorReasonMISSING_DESCRIPTION2 AdErrorReason = "MISSING_DESCRIPTION2"

	//
	// The destination url must contain at least one tag (e.g. {lpurl})
	//
	AdErrorReasonMISSING_DESTINATION_URL_TAG AdErrorReason = "MISSING_DESTINATION_URL_TAG"

	//
	// The tracking url template of ExpandedDynamicSearchAd must contain at least one tag.
	// (e.g. {lpurl})
	//
	AdErrorReasonMISSING_LANDING_PAGE_URL_TAG AdErrorReason = "MISSING_LANDING_PAGE_URL_TAG"

	//
	// A valid dimension must be specified for this ad.
	//
	AdErrorReasonMISSING_DIMENSION AdErrorReason = "MISSING_DIMENSION"

	//
	// A display URL must be entered.
	//
	AdErrorReasonMISSING_DISPLAY_URL AdErrorReason = "MISSING_DISPLAY_URL"

	//
	// Headline must be entered.
	//
	AdErrorReasonMISSING_HEADLINE AdErrorReason = "MISSING_HEADLINE"

	//
	// A height must be entered.
	//
	AdErrorReasonMISSING_HEIGHT AdErrorReason = "MISSING_HEIGHT"

	//
	// An image must be entered.
	//
	AdErrorReasonMISSING_IMAGE AdErrorReason = "MISSING_IMAGE"

	//
	// Marketing image or product videos are required.
	//
	AdErrorReasonMISSING_MARKETING_IMAGE_OR_PRODUCT_VIDEOS AdErrorReason = "MISSING_MARKETING_IMAGE_OR_PRODUCT_VIDEOS"

	//
	// The markup language in which your site is written must be entered.
	//
	AdErrorReasonMISSING_MARKUP_LANGUAGES AdErrorReason = "MISSING_MARKUP_LANGUAGES"

	//
	// A mobile carrier must be entered.
	//
	AdErrorReasonMISSING_MOBILE_CARRIER AdErrorReason = "MISSING_MOBILE_CARRIER"

	//
	// Phone number must be entered.
	//
	AdErrorReasonMISSING_PHONE AdErrorReason = "MISSING_PHONE"

	//
	// Missing required template fields
	//
	AdErrorReasonMISSING_REQUIRED_TEMPLATE_FIELDS AdErrorReason = "MISSING_REQUIRED_TEMPLATE_FIELDS"

	//
	// Missing a required field value
	//
	AdErrorReasonMISSING_TEMPLATE_FIELD_VALUE AdErrorReason = "MISSING_TEMPLATE_FIELD_VALUE"

	//
	// The ad must have text.
	//
	AdErrorReasonMISSING_TEXT AdErrorReason = "MISSING_TEXT"

	//
	// A visible URL must be entered.
	//
	AdErrorReasonMISSING_VISIBLE_URL AdErrorReason = "MISSING_VISIBLE_URL"

	//
	// A width must be entered.
	//
	AdErrorReasonMISSING_WIDTH AdErrorReason = "MISSING_WIDTH"

	//
	// Only 1 feed can be used as the source of ad customizer substitutions in a single ad.
	//
	AdErrorReasonMULTIPLE_DISTINCT_FEEDS_UNSUPPORTED AdErrorReason = "MULTIPLE_DISTINCT_FEEDS_UNSUPPORTED"

	//
	// TempAdUnionId must be use when adding template ads.
	//
	AdErrorReasonMUST_USE_TEMP_AD_UNION_ID_ON_ADD AdErrorReason = "MUST_USE_TEMP_AD_UNION_ID_ON_ADD"

	//
	// The string has too many characters.
	//
	AdErrorReasonTOO_LONG AdErrorReason = "TOO_LONG"

	//
	// The string has too few characters.
	//
	AdErrorReasonTOO_SHORT AdErrorReason = "TOO_SHORT"

	//
	// Ad union dimensions cannot change for saved ads.
	//
	AdErrorReasonUNION_DIMENSIONS_CANNOT_CHANGE AdErrorReason = "UNION_DIMENSIONS_CANNOT_CHANGE"

	//
	// Address component is not {country, lat, lng}.
	//
	AdErrorReasonUNKNOWN_ADDRESS_COMPONENT AdErrorReason = "UNKNOWN_ADDRESS_COMPONENT"

	//
	// Unknown unique field name
	//
	AdErrorReasonUNKNOWN_FIELD_NAME AdErrorReason = "UNKNOWN_FIELD_NAME"

	//
	// Unknown unique name (template element type specifier)
	//
	AdErrorReasonUNKNOWN_UNIQUE_NAME AdErrorReason = "UNKNOWN_UNIQUE_NAME"

	//
	// Unsupported ad dimension
	//
	AdErrorReasonUNSUPPORTED_DIMENSIONS AdErrorReason = "UNSUPPORTED_DIMENSIONS"

	//
	// URL starts with an invalid scheme.
	//
	AdErrorReasonURL_INVALID_SCHEME AdErrorReason = "URL_INVALID_SCHEME"

	//
	// URL ends with an invalid top-level domain name.
	//
	AdErrorReasonURL_INVALID_TOP_LEVEL_DOMAIN AdErrorReason = "URL_INVALID_TOP_LEVEL_DOMAIN"

	//
	// URL contains illegal characters.
	//
	AdErrorReasonURL_MALFORMED AdErrorReason = "URL_MALFORMED"

	//
	// URL must contain a host name.
	//
	AdErrorReasonURL_NO_HOST AdErrorReason = "URL_NO_HOST"

	//
	// URL not equivalent during upgrade.
	//
	AdErrorReasonURL_NOT_EQUIVALENT AdErrorReason = "URL_NOT_EQUIVALENT"

	//
	// URL host name too long to be stored as visible URL (applies to Ad Exchange ads)
	//
	AdErrorReasonURL_HOST_NAME_TOO_LONG AdErrorReason = "URL_HOST_NAME_TOO_LONG"

	//
	// URL must start with a scheme.
	//
	AdErrorReasonURL_NO_SCHEME AdErrorReason = "URL_NO_SCHEME"

	//
	// URL should end in a valid domain extension, such as .com or .net.
	//
	AdErrorReasonURL_NO_TOP_LEVEL_DOMAIN AdErrorReason = "URL_NO_TOP_LEVEL_DOMAIN"

	//
	// URL must not end with a path.
	//
	AdErrorReasonURL_PATH_NOT_ALLOWED AdErrorReason = "URL_PATH_NOT_ALLOWED"

	//
	// URL must not specify a port.
	//
	AdErrorReasonURL_PORT_NOT_ALLOWED AdErrorReason = "URL_PORT_NOT_ALLOWED"

	//
	// URL must not contain a query.
	//
	AdErrorReasonURL_QUERY_NOT_ALLOWED AdErrorReason = "URL_QUERY_NOT_ALLOWED"

	//
	// A url scheme is not allowed in front of tag in dest url (e.g. http://{lpurl})
	//
	AdErrorReasonURL_SCHEME_BEFORE_DSA_TAG AdErrorReason = "URL_SCHEME_BEFORE_DSA_TAG"

	//
	// A url scheme is not allowed in front of tag in tracking url template (e.g. http://{lpurl})
	//
	AdErrorReasonURL_SCHEME_BEFORE_EXPANDED_DYNAMIC_SEARCH_AD_TAG AdErrorReason = "URL_SCHEME_BEFORE_EXPANDED_DYNAMIC_SEARCH_AD_TAG"

	//
	// The user does not have permissions to create a template ad for the given
	// template.
	//
	AdErrorReasonUSER_DOES_NOT_HAVE_ACCESS_TO_TEMPLATE AdErrorReason = "USER_DOES_NOT_HAVE_ACCESS_TO_TEMPLATE"

	//
	// Expandable setting is inconsistent/wrong. For example, an AdX ad is
	// invalid if it has a expandable vendor format but no expanding directions
	// specified, or expanding directions is specified, but the vendor format
	// is not expandable.
	//
	AdErrorReasonINCONSISTENT_EXPANDABLE_SETTINGS AdErrorReason = "INCONSISTENT_EXPANDABLE_SETTINGS"

	//
	// Format is invalid
	//
	AdErrorReasonINVALID_FORMAT AdErrorReason = "INVALID_FORMAT"

	//
	// The text of this field did not match a pattern of allowed values.
	//
	AdErrorReasonINVALID_FIELD_TEXT AdErrorReason = "INVALID_FIELD_TEXT"

	//
	// Template element is mising
	//
	AdErrorReasonELEMENT_NOT_PRESENT AdErrorReason = "ELEMENT_NOT_PRESENT"

	//
	// Error occurred during image processing
	//
	AdErrorReasonIMAGE_ERROR AdErrorReason = "IMAGE_ERROR"

	//
	// The value is not within the valid range
	//
	AdErrorReasonVALUE_NOT_IN_RANGE AdErrorReason = "VALUE_NOT_IN_RANGE"

	//
	// Template element field is not present
	//
	AdErrorReasonFIELD_NOT_PRESENT AdErrorReason = "FIELD_NOT_PRESENT"

	//
	// Address is incomplete
	//
	AdErrorReasonADDRESS_NOT_COMPLETE AdErrorReason = "ADDRESS_NOT_COMPLETE"

	//
	// Invalid address
	//
	AdErrorReasonADDRESS_INVALID AdErrorReason = "ADDRESS_INVALID"

	//
	// Error retrieving specified video
	//
	AdErrorReasonVIDEO_RETRIEVAL_ERROR AdErrorReason = "VIDEO_RETRIEVAL_ERROR"

	//
	// Error processing audio
	//
	AdErrorReasonAUDIO_ERROR AdErrorReason = "AUDIO_ERROR"

	//
	// Display URL is incorrect for YouTube PYV ads
	//
	AdErrorReasonINVALID_YOUTUBE_DISPLAY_URL AdErrorReason = "INVALID_YOUTUBE_DISPLAY_URL"

	//
	// Too many product Images in GmailAd
	//
	AdErrorReasonTOO_MANY_PRODUCT_IMAGES AdErrorReason = "TOO_MANY_PRODUCT_IMAGES"

	//
	// Too many product Videos in GmailAd
	//
	AdErrorReasonTOO_MANY_PRODUCT_VIDEOS AdErrorReason = "TOO_MANY_PRODUCT_VIDEOS"

	//
	// The device preference is not compatible with the ad type
	//
	AdErrorReasonINCOMPATIBLE_AD_TYPE_AND_DEVICE_PREFERENCE AdErrorReason = "INCOMPATIBLE_AD_TYPE_AND_DEVICE_PREFERENCE"

	//
	// Call tracking is not supported for specified country.
	//
	AdErrorReasonCALLTRACKING_NOT_SUPPORTED_FOR_COUNTRY AdErrorReason = "CALLTRACKING_NOT_SUPPORTED_FOR_COUNTRY"

	//
	// Carrier specific short number is not allowed.
	//
	AdErrorReasonCARRIER_SPECIFIC_SHORT_NUMBER_NOT_ALLOWED AdErrorReason = "CARRIER_SPECIFIC_SHORT_NUMBER_NOT_ALLOWED"

	//
	// Specified phone number type is disallowed.
	//
	AdErrorReasonDISALLOWED_NUMBER_TYPE AdErrorReason = "DISALLOWED_NUMBER_TYPE"

	//
	// Phone number not supported for country.
	//
	AdErrorReasonPHONE_NUMBER_NOT_SUPPORTED_FOR_COUNTRY AdErrorReason = "PHONE_NUMBER_NOT_SUPPORTED_FOR_COUNTRY"

	//
	// Phone number not supported with call tracking enabled for country.
	//
	AdErrorReasonPHONE_NUMBER_NOT_SUPPORTED_WITH_CALLTRACKING_FOR_COUNTRY AdErrorReason = "PHONE_NUMBER_NOT_SUPPORTED_WITH_CALLTRACKING_FOR_COUNTRY"

	//
	// Premium rate phone number is not allowed.
	//
	AdErrorReasonPREMIUM_RATE_NUMBER_NOT_ALLOWED AdErrorReason = "PREMIUM_RATE_NUMBER_NOT_ALLOWED"

	//
	// Vanity phone number is not allowed.
	//
	AdErrorReasonVANITY_PHONE_NUMBER_NOT_ALLOWED AdErrorReason = "VANITY_PHONE_NUMBER_NOT_ALLOWED"

	//
	// Invalid call conversion type id.
	//
	AdErrorReasonINVALID_CALL_CONVERSION_TYPE_ID AdErrorReason = "INVALID_CALL_CONVERSION_TYPE_ID"

	AdErrorReasonCANNOT_DISABLE_CALL_CONVERSION_AND_SET_CONVERSION_TYPE_ID AdErrorReason = "CANNOT_DISABLE_CALL_CONVERSION_AND_SET_CONVERSION_TYPE_ID"

	//
	// Cannot set path2 without path1.
	//
	AdErrorReasonCANNOT_SET_PATH2_WITHOUT_PATH1 AdErrorReason = "CANNOT_SET_PATH2_WITHOUT_PATH1"

	//
	// Missing domain name in campaign setting when adding expanded dynamic search ad.
	//
	AdErrorReasonMISSING_DYNAMIC_SEARCH_ADS_SETTING_DOMAIN_NAME AdErrorReason = "MISSING_DYNAMIC_SEARCH_ADS_SETTING_DOMAIN_NAME"

	//
	// An unexpected or unknown error occurred.
	//
	AdErrorReasonUNKNOWN AdErrorReason = "UNKNOWN"
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
// The reasons for the target error.
//
type CollectionSizeErrorReason string

const (
	CollectionSizeErrorReasonTOO_FEW CollectionSizeErrorReason = "TOO_FEW"

	CollectionSizeErrorReasonTOO_MANY CollectionSizeErrorReason = "TOO_MANY"
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

type EntityNotFoundReason string

const (

	//
	// The specified id refered to an entity which either doesn't exist or is not accessible to the
	// customer. e.g. campaign belongs to another customer.
	//
	EntityNotFoundReasonINVALID_ID EntityNotFoundReason = "INVALID_ID"
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
// The single reason for the whitelist error.
//
type NotWhitelistedErrorReason string

const (

	//
	// Customer is not whitelisted for accessing the API.
	//
	NotWhitelistedErrorReasonCUSTOMER_NOT_WHITELISTED_FOR_API NotWhitelistedErrorReason = "CUSTOMER_NOT_WHITELISTED_FOR_API"
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

type AdError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AdErrorReason `xml:"reason,omitempty"`
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

type CollectionSizeError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CollectionSizeError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *CollectionSizeErrorReason `xml:"reason,omitempty"`
}

type ComparableValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ComparableValue"`

	//
	// Indicates that this instance is a subtype of ComparableValue.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	ComparableValueType string `xml:"ComparableValue.Type,omitempty"`
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

type DoubleValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DoubleValue"`

	*NumberValue

	//
	// the underlying double value.
	//
	Number float64 `xml:"number,omitempty"`
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

type LongValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 LongValue"`

	*NumberValue

	//
	// the underlying long value.
	//
	Number int64 `xml:"number,omitempty"`
}

type Money struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Money"`

	*ComparableValue

	//
	// Amount in micros. One million is equivalent to one unit.
	//
	MicroAmount int64 `xml:"microAmount,omitempty"`
}

type NotEmptyError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NotEmptyError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NotEmptyErrorReason `xml:"reason,omitempty"`
}

type NotWhitelistedError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NotWhitelistedError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NotWhitelistedErrorReason `xml:"reason,omitempty"`
}

type NullError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NullError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *NullErrorReason `xml:"reason,omitempty"`
}

type NumberValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 NumberValue"`

	*ComparableValue
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

//
// Encodes the reason (cause) of a particular {@link CurrencyCodeError}.
//
type CurrencyCodeErrorReason string

const (
	CurrencyCodeErrorReasonUNSUPPORTED_CURRENCY_CODE CurrencyCodeErrorReason = "UNSUPPORTED_CURRENCY_CODE"
)

//
// Reasons
//
type OfflineDataUploadErrorReason string

const (
	OfflineDataUploadErrorReasonUNKNOWN OfflineDataUploadErrorReason = "UNKNOWN"

	//
	// Indicates a row error due to the incompatible {@code OfflineDataUploadUserIdentifierType},
	// like using EXTERNAL_USER_ID for first party uploads or not using EXTERNAL_USER_ID for third
	// party uploads.
	//
	OfflineDataUploadErrorReasonINCOMPATIBLE_USERIDENTIFIER_TYPE OfflineDataUploadErrorReason = "INCOMPATIBLE_USERIDENTIFIER_TYPE"

	//
	// Indicates an upload error due to the invalid upload type.
	//
	OfflineDataUploadErrorReasonINVALID_UPLOAD_TYPE OfflineDataUploadErrorReason = "INVALID_UPLOAD_TYPE"

	//
	// Indicates an upload error due to missing metadata.
	//
	OfflineDataUploadErrorReasonMISSING_UPLOAD_METADATA OfflineDataUploadErrorReason = "MISSING_UPLOAD_METADATA"

	//
	// Indicates an upload error due to missing metadata.
	//
	OfflineDataUploadErrorReasonINVALID_UPLOAD_METADATA OfflineDataUploadErrorReason = "INVALID_UPLOAD_METADATA"

	//
	// Indicates an upload error due to invalid partner id in metadata.
	//
	OfflineDataUploadErrorReasonINVALID_PARTNER_ID OfflineDataUploadErrorReason = "INVALID_PARTNER_ID"

	//
	// Indicates a row error due to missing transaction data.
	//
	OfflineDataUploadErrorReasonMISSING_TRANSACTION_INFO OfflineDataUploadErrorReason = "MISSING_TRANSACTION_INFO"

	//
	// The name specified in store_sales_attributes is used to report conversions to a conversion
	// type configured in AdWords with the same name. A row generates this error if there is no such
	// name configured in the account.
	//
	OfflineDataUploadErrorReasonINVALID_CONVERSION_TYPE OfflineDataUploadErrorReason = "INVALID_CONVERSION_TYPE"

	//
	// Indicates a row error due to a conversion with a transaction time in the future.
	//
	OfflineDataUploadErrorReasonFUTURE_TRANSACTION_TIME OfflineDataUploadErrorReason = "FUTURE_TRANSACTION_TIME"

	//
	// Indicates a row error due to a negative transaction amount.
	//
	OfflineDataUploadErrorReasonNEGATIVE_TRANSACTION_AMOUNT OfflineDataUploadErrorReason = "NEGATIVE_TRANSACTION_AMOUNT"

	//
	// Country code hashed.
	//
	OfflineDataUploadErrorReasonCOUNTRY_CODE_HASHED OfflineDataUploadErrorReason = "COUNTRY_CODE_HASHED"

	//
	// ZIP Code hashed.
	//
	OfflineDataUploadErrorReasonZIPCODE_HASHED OfflineDataUploadErrorReason = "ZIPCODE_HASHED"

	//
	// Email not hashed.
	//
	OfflineDataUploadErrorReasonEMAIL_NOT_HASHED OfflineDataUploadErrorReason = "EMAIL_NOT_HASHED"

	//
	// First Name not hashed.
	//
	OfflineDataUploadErrorReasonFIRST_NAME_NOT_HASHED OfflineDataUploadErrorReason = "FIRST_NAME_NOT_HASHED"

	//
	// Last Name not hashed.
	//
	OfflineDataUploadErrorReasonLAST_NAME_NOT_HASHED OfflineDataUploadErrorReason = "LAST_NAME_NOT_HASHED"

	//
	// Phone not hashed.
	//
	OfflineDataUploadErrorReasonPHONE_NOT_HASHED OfflineDataUploadErrorReason = "PHONE_NOT_HASHED"
)

//
// Indicates the offline data upload processing failure reason.
//
type OfflineDataUploadFailureReason string

const (

	//
	// UNKNOWN value cannot be passed as input.
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	OfflineDataUploadFailureReasonUNKNOWN OfflineDataUploadFailureReason = "UNKNOWN"

	//
	// Indicates the matched transactions don?t cross the minimum threshold.
	//
	OfflineDataUploadFailureReasonINSUFFICIENT_MATCHED_TRANSACTIONS OfflineDataUploadFailureReason = "INSUFFICIENT_MATCHED_TRANSACTIONS"

	//
	// Indicates the insufficient transactions uploaded.
	//
	OfflineDataUploadFailureReasonINSUFFICIENT_TRANSACTIONS OfflineDataUploadFailureReason = "INSUFFICIENT_TRANSACTIONS"
)

//
// This indicates the status of offline upload.
//
type OfflineDataUploadStatus string

const (

	//
	// UNKNOWN value cannot be passed as input.
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	OfflineDataUploadStatusUNKNOWN OfflineDataUploadStatus = "UNKNOWN"

	//
	// Indicates the upload failed in the offline processing.
	//
	OfflineDataUploadStatusFAILURE OfflineDataUploadStatus = "FAILURE"

	//
	// Indicates the upload passed formatting checks and was accepted for offline
	// processing.
	//
	OfflineDataUploadStatusIN_PROCESS OfflineDataUploadStatus = "IN_PROCESS"

	//
	// Indicates the upload was processed by the offline processing pipeline.
	//
	OfflineDataUploadStatusSUCCESS OfflineDataUploadStatus = "SUCCESS"
)

//
// Upload types.
//
type OfflineDataUploadType string

const (

	//
	// UNKNOWN value cannot be passed as input.
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	OfflineDataUploadTypeUNKNOWN OfflineDataUploadType = "UNKNOWN"

	//
	// Indicates Store Sales Direct Upload for self service.
	//
	OfflineDataUploadTypeSTORE_SALES_UPLOAD_FIRST_PARTY OfflineDataUploadType = "STORE_SALES_UPLOAD_FIRST_PARTY"

	//
	// Indicates Store Sales Direct Upload for third party.
	//
	OfflineDataUploadTypeSTORE_SALES_UPLOAD_THIRD_PARTY OfflineDataUploadType = "STORE_SALES_UPLOAD_THIRD_PARTY"
)

//
// Indentifier types of user information.
//
type OfflineDataUploadUserIdentifierType string

const (

	//
	// UNKNOWN value can not be passed as input.
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	OfflineDataUploadUserIdentifierTypeUNKNOWN OfflineDataUploadUserIdentifierType = "UNKNOWN"

	//
	// Indicates the email address.
	//
	OfflineDataUploadUserIdentifierTypeHASHED_EMAIL OfflineDataUploadUserIdentifierType = "HASHED_EMAIL"

	//
	// Indicates the phone number.
	//
	OfflineDataUploadUserIdentifierTypeHASHED_PHONE OfflineDataUploadUserIdentifierType = "HASHED_PHONE"

	//
	// Indicates the last name.
	//
	OfflineDataUploadUserIdentifierTypeHASHED_LAST_NAME OfflineDataUploadUserIdentifierType = "HASHED_LAST_NAME"

	//
	// Indicates the first name.
	//
	OfflineDataUploadUserIdentifierTypeHASHED_FIRST_NAME OfflineDataUploadUserIdentifierType = "HASHED_FIRST_NAME"

	//
	// Indicates the city.
	//
	OfflineDataUploadUserIdentifierTypeCITY OfflineDataUploadUserIdentifierType = "CITY"

	//
	// Indicates the state.
	//
	OfflineDataUploadUserIdentifierTypeSTATE OfflineDataUploadUserIdentifierType = "STATE"

	//
	// Indicates the zip code.
	//
	OfflineDataUploadUserIdentifierTypeZIPCODE OfflineDataUploadUserIdentifierType = "ZIPCODE"

	//
	// ISO two-letter country codes.
	//
	OfflineDataUploadUserIdentifierTypeCOUNTRY_CODE OfflineDataUploadUserIdentifierType = "COUNTRY_CODE"

	//
	// Indicates the external id like third party id.
	//
	OfflineDataUploadUserIdentifierTypeEXTERNAL_USER_ID OfflineDataUploadUserIdentifierType = "EXTERNAL_USER_ID"
)

type Get struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 get"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	ServiceSelector *Selector `xml:"serviceSelector,omitempty"`
}

type GetResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 getResponse"`

	Rval *OfflineDataUploadPage `xml:"rval,omitempty"`
}

type Mutate struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 mutate"`

	//
	// <span class="constraint CollectionSize">The minimum size of this collection is 1.</span>
	// <span class="constraint ContentsNotNull">This field must not contain {@code null} elements.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	// <span class="constraint SupportedOperators">The following {@link Operator}s are supported: ADD, SET.</span>
	//
	Operations []*OfflineDataUploadOperation `xml:"operations,omitempty"`
}

type MutateResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 mutateResponse"`

	Rval *OfflineDataUploadReturnValue `xml:"rval,omitempty"`
}

type CurrencyCodeError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 CurrencyCodeError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *CurrencyCodeErrorReason `xml:"reason,omitempty"`
}

type FirstPartyUploadMetadata struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 FirstPartyUploadMetadata"`

	*StoreSalesUploadCommonMetadata
}

type MoneyWithCurrency struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 MoneyWithCurrency"`

	*ComparableValue

	//
	// The amount of money.
	// <span class="constraint InRange">This field must be greater than or equal to 0.</span>
	//
	Money *Money `xml:"money,omitempty"`

	//
	// Currency code.
	// <span class="constraint StringLength">The length of this string should be between 3 and 3, inclusive, (trimmed).</span>
	//
	CurrencyCode string `xml:"currencyCode,omitempty"`
}

type OfflineData struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 OfflineData"`

	StoreSalesTransaction *StoreSalesTransaction `xml:"StoreSalesTransaction,omitempty"`
}

type OfflineDataUpload struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 OfflineDataUpload"`

	//
	// User specified upload id.
	// <span class="constraint Selectable">This field can be selected using the value "ExternalUploadId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	ExternalUploadId int64 `xml:"externalUploadId,omitempty"`

	//
	// Type of this upload.
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	UploadType *OfflineDataUploadType `xml:"uploadType,omitempty"`

	//
	// Status of this upload.
	// <span class="constraint Selectable">This field can be selected using the value "UploadStatus".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: ADD.</span>
	//
	UploadStatus *OfflineDataUploadStatus `xml:"uploadStatus,omitempty"`

	//
	// Metadata for this upload.
	//
	UploadMetadata *UploadMetadata `xml:"uploadMetadata,omitempty"`

	//
	// List of offline data in this upload. For AdWords API, each offlineDataList can have at most 50
	// OfflineData.
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : ADD.</span>
	//
	OfflineDataList []*OfflineData `xml:"offlineDataList,omitempty"`

	//
	// Processing failure reason for get, if status is FAILURE. Used for upload level failures.
	// <span class="constraint Selectable">This field can be selected using the value "FailureReason".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: ADD.</span>
	//
	FailureReason *OfflineDataUploadFailureReason `xml:"failureReason,omitempty"`
}

type OfflineDataUploadError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 OfflineDataUploadError"`

	*ApiError

	Reason *OfflineDataUploadErrorReason `xml:"reason,omitempty"`
}

type OfflineDataUploadOperation struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 OfflineDataUploadOperation"`

	*Operation

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Operand *OfflineDataUpload `xml:"operand,omitempty"`
}

type OfflineDataUploadPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 OfflineDataUploadPage"`

	*Page

	Entries []*OfflineDataUpload `xml:"entries,omitempty"`
}

type OfflineDataUploadReturnValue struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 OfflineDataUploadReturnValue"`

	*ListReturnValue

	Value []*OfflineDataUpload `xml:"value,omitempty"`

	//
	// In v201710 and previous, this field stores a list of operation level errors. Starting in
	// v201802, this field stores both operation-level and row-level errors. For row-level errors,
	// offlineDataList will be shown in the fieldPath along with row index. In this case, the
	// operation will be processed and just the rows with errors will not be used. For more
	// information about partial failure, see:
	// https://developers.google.com/adwords/api/docs/guides/partial-failure
	//
	PartialFailureErrors []*ApiError `xml:"partialFailureErrors,omitempty"`
}

type StoreSalesTransaction struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 StoreSalesTransaction"`

	//
	// List of UserIdentifiers.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	UserIdentifiers []*UserIdentifier `xml:"userIdentifiers,omitempty"`

	//
	// Transaction time.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	TransactionTime string `xml:"transactionTime,omitempty"`

	//
	// Transaction amount. We support the ISO 4217 3-character currency code. For example: USD, EUR.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	TransactionAmount *MoneyWithCurrency `xml:"transactionAmount,omitempty"`

	//
	// Conversion name configured while creating ConversionType in AdWords.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	ConversionName string `xml:"conversionName,omitempty"`
}

type StoreSalesUploadCommonMetadata struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 StoreSalesUploadCommonMetadata"`

	//
	// This is the fraction of overall sales which you can associate with a customer loyalty program.
	// The fraction needs to be between 0 and 1 (excluding 0).
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	LoyaltyRate float64 `xml:"loyaltyRate,omitempty"`

	//
	// This is the ratio of sales you?re uploading compared to the overall sales that you can
	// associate with a customer. The fraction needs to be between 0 and 1. For example, if you upload
	// half the sales that you are able to associate with a customer, your Transaction Upload Rate
	// would be 0.5 (excluding 0).
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	TransactionUploadRate float64 `xml:"transactionUploadRate,omitempty"`

	//
	// Indicates that this instance is a subtype of StoreSalesUploadCommonMetadata.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	StoreSalesUploadCommonMetadataType string `xml:"StoreSalesUploadCommonMetadata.Type,omitempty"`
}

type ThirdPartyUploadMetadata struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 ThirdPartyUploadMetadata"`

	*StoreSalesUploadCommonMetadata

	//
	// Advertiser upload time to partner.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	AdvertiserUploadTime string `xml:"advertiserUploadTime,omitempty"`

	//
	// The fraction of transactions that are valid. Invalid transactions may include invalid format,
	// values. Range (0.0 to 1.0]
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	ValidTransactionRate float64 `xml:"validTransactionRate,omitempty"`

	//
	// The fraction of valid transactions that are matched to an external user id on the partner side.
	// Range (0.0 to 1.0]
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	PartnerMatchRate float64 `xml:"partnerMatchRate,omitempty"`

	//
	// The fraction of valid transactions that are uploaded by the partner to Google. Range (0.0 to
	// 1.0]
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	PartnerUploadRate float64 `xml:"partnerUploadRate,omitempty"`

	//
	// Indicates the version of partnerIds to be used for uploads.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	BridgeMapVersionId string `xml:"bridgeMapVersionId,omitempty"`

	//
	// The ID of the third party partner uploading the transaction feed.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	PartnerId int32 `xml:"partnerId,omitempty"`
}

type UploadMetadata struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 UploadMetadata"`

	StoreSalesUploadCommonMetadata *StoreSalesUploadCommonMetadata `xml:"StoreSalesUploadCommonMetadata,omitempty"`
}

type UserIdentifier struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/rm/v201802 UserIdentifier"`

	//
	// Type of user identifier.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	UserIdentifierType *OfflineDataUploadUserIdentifierType `xml:"userIdentifierType,omitempty"`

	//
	// Value of identifier. Hashed using SHA-256 if needed.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Value string `xml:"value,omitempty"`
}

type OfflineDataUploadServiceInterface struct {
	client *SOAPClient
}

func NewOfflineDataUploadServiceInterface(url string, tls bool, auth *BasicAuth) *OfflineDataUploadServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &OfflineDataUploadServiceInterface{
		client: client,
	}
}

func NewOfflineDataUploadServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *OfflineDataUploadServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &OfflineDataUploadServiceInterface{
		client: client,
	}
}

func (service *OfflineDataUploadServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *OfflineDataUploadServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Returns a list of OfflineDataUpload objects that match the criteria specified in the selector.

   <p><b>Note:</b> If an upload fails after processing, reason will be reported in {@link
   OfflineDataUpload#failureReason}.

   @throws {@link ApiException} if problems occurred while retrieving results.
*/
func (service *OfflineDataUploadServiceInterface) Get(request *Get) (*GetResponse, error) {
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
   Applies a list of mutate operations (i.e. add, set) to offline data upload:

   <p>Add - uploads offline data for each entry in operations. Some operations can fail for upload
   level errors like invalid {@code UploadMetadata}. Check {@code OfflineDataUploadReturnValue}
   for partial failure list.

   <p>Set - updates the upload result for each upload. It is for internal use only.

   <p><b>Note:</b> For AdWords API, one ADD request can have at most 2000 operations.

   <p><b>Note:</b> Add operation might possibly succeed even with errors in {@code OfflineData}.
   Data errors are reported in {@link OfflineDataUpload#partialDataErrors}

   <p><b>Note:</b> Supports only the {@code ADD} operator. {@code SET} operator is internally used
   only.({@code REMOVE} is not supported).

   @param operations A list of offline data upload operations.
   @return The list of offline data upload results in the same order as operations.
   @throws {@link ApiException} if problems occur.
*/
func (service *OfflineDataUploadServiceInterface) Mutate(request *Mutate) (*MutateResponse, error) {
	response := new(MutateResponse)
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
