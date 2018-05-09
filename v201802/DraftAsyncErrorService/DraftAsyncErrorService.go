package DraftAsyncErrorService

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
// The reasons for the target error.
//
type AdGroupAdErrorReason string

const (

	//
	// No link found between the adgroup ad and the label.
	//
	AdGroupAdErrorReasonAD_GROUP_AD_LABEL_DOES_NOT_EXIST AdGroupAdErrorReason = "AD_GROUP_AD_LABEL_DOES_NOT_EXIST"

	//
	// The label has already been attached to the adgroup ad.
	//
	AdGroupAdErrorReasonAD_GROUP_AD_LABEL_ALREADY_EXISTS AdGroupAdErrorReason = "AD_GROUP_AD_LABEL_ALREADY_EXISTS"

	//
	// The specified ad was not found in the adgroup
	//
	AdGroupAdErrorReasonAD_NOT_UNDER_ADGROUP AdGroupAdErrorReason = "AD_NOT_UNDER_ADGROUP"

	//
	// Removed ads may not be modified
	//
	AdGroupAdErrorReasonCANNOT_OPERATE_ON_REMOVED_ADGROUPAD AdGroupAdErrorReason = "CANNOT_OPERATE_ON_REMOVED_ADGROUPAD"

	//
	// An ad of this type is deprecated and cannot be created. Only deletions
	// are permitted.
	//
	AdGroupAdErrorReasonCANNOT_CREATE_DEPRECATED_ADS AdGroupAdErrorReason = "CANNOT_CREATE_DEPRECATED_ADS"

	//
	// Text ads are deprecated and cannot be created. Use expanded text ads instead.
	//
	AdGroupAdErrorReasonCANNOT_CREATE_TEXT_ADS AdGroupAdErrorReason = "CANNOT_CREATE_TEXT_ADS"

	//
	// A required field was not specified or is an empty string.
	//
	AdGroupAdErrorReasonEMPTY_FIELD AdGroupAdErrorReason = "EMPTY_FIELD"

	//
	// An ad may only be modified once per call
	//
	AdGroupAdErrorReasonENTITY_REFERENCED_IN_MULTIPLE_OPS AdGroupAdErrorReason = "ENTITY_REFERENCED_IN_MULTIPLE_OPS"

	//
	// The specified operation is not supported.  Only ADD, SET, and REMOVE
	// are supported
	//
	AdGroupAdErrorReasonUNSUPPORTED_OPERATION AdGroupAdErrorReason = "UNSUPPORTED_OPERATION"
)

//
// The reasons for the target error.
//
type AdGroupCriterionErrorReason string

const (

	//
	// No link found between the AdGroupCriterion and the label.
	//
	AdGroupCriterionErrorReasonAD_GROUP_CRITERION_LABEL_DOES_NOT_EXIST AdGroupCriterionErrorReason = "AD_GROUP_CRITERION_LABEL_DOES_NOT_EXIST"

	//
	// The label has already been attached to the AdGroupCriterion.
	//
	AdGroupCriterionErrorReasonAD_GROUP_CRITERION_LABEL_ALREADY_EXISTS AdGroupCriterionErrorReason = "AD_GROUP_CRITERION_LABEL_ALREADY_EXISTS"

	//
	// Negative AdGroupCriterion cannot have labels.
	//
	AdGroupCriterionErrorReasonCANNOT_ADD_LABEL_TO_NEGATIVE_CRITERION AdGroupCriterionErrorReason = "CANNOT_ADD_LABEL_TO_NEGATIVE_CRITERION"

	//
	// Too many operations for a single call.
	//
	AdGroupCriterionErrorReasonTOO_MANY_OPERATIONS AdGroupCriterionErrorReason = "TOO_MANY_OPERATIONS"

	//
	// Negative ad group criteria are not updateable.
	//
	AdGroupCriterionErrorReasonCANT_UPDATE_NEGATIVE AdGroupCriterionErrorReason = "CANT_UPDATE_NEGATIVE"

	//
	// Concrete type of criterion (keyword v.s. placement) is required for
	// ADD and SET operations.
	//
	AdGroupCriterionErrorReasonCONCRETE_TYPE_REQUIRED AdGroupCriterionErrorReason = "CONCRETE_TYPE_REQUIRED"

	//
	// Bid is incompatible with ad group's bidding settings.
	//
	AdGroupCriterionErrorReasonBID_INCOMPATIBLE_WITH_ADGROUP AdGroupCriterionErrorReason = "BID_INCOMPATIBLE_WITH_ADGROUP"

	//
	// Cannot target and exclude the same criterion at once.
	//
	AdGroupCriterionErrorReasonCANNOT_TARGET_AND_EXCLUDE AdGroupCriterionErrorReason = "CANNOT_TARGET_AND_EXCLUDE"

	//
	// The URL of a placement is invalid.
	//
	AdGroupCriterionErrorReasonILLEGAL_URL AdGroupCriterionErrorReason = "ILLEGAL_URL"

	//
	// Keyword text was invalid.
	//
	AdGroupCriterionErrorReasonINVALID_KEYWORD_TEXT AdGroupCriterionErrorReason = "INVALID_KEYWORD_TEXT"

	//
	// Destination URL was invalid.
	//
	AdGroupCriterionErrorReasonINVALID_DESTINATION_URL AdGroupCriterionErrorReason = "INVALID_DESTINATION_URL"

	//
	// The destination url must contain at least one tag (e.g. {lpurl})
	//
	AdGroupCriterionErrorReasonMISSING_DESTINATION_URL_TAG AdGroupCriterionErrorReason = "MISSING_DESTINATION_URL_TAG"

	//
	// Keyword-level cpm bid is not supported
	//
	AdGroupCriterionErrorReasonKEYWORD_LEVEL_BID_NOT_SUPPORTED_FOR_MANUALCPM AdGroupCriterionErrorReason = "KEYWORD_LEVEL_BID_NOT_SUPPORTED_FOR_MANUALCPM"

	//
	// For example, cannot add a biddable ad group criterion that had been removed.
	//
	AdGroupCriterionErrorReasonINVALID_USER_STATUS AdGroupCriterionErrorReason = "INVALID_USER_STATUS"

	//
	// Criteria type cannot be targeted for the ad group. Either the account
	// is restricted to keywords only, the criteria type is incompatible
	// with the campaign's bidding strategy, or the criteria type can only
	// be applied to campaigns.
	//
	AdGroupCriterionErrorReasonCANNOT_ADD_CRITERIA_TYPE AdGroupCriterionErrorReason = "CANNOT_ADD_CRITERIA_TYPE"

	//
	// Criteria type cannot be excluded for the ad group. Refer to the
	// documentation for a specific criterion to check if it is excludable.
	//
	AdGroupCriterionErrorReasonCANNOT_EXCLUDE_CRITERIA_TYPE AdGroupCriterionErrorReason = "CANNOT_EXCLUDE_CRITERIA_TYPE"

	//
	// Ad group is invalid due to the product partitions it contains.
	//
	AdGroupCriterionErrorReasonINVALID_PRODUCT_PARTITION_HIERARCHY AdGroupCriterionErrorReason = "INVALID_PRODUCT_PARTITION_HIERARCHY"

	//
	// Product partition unit cannot have children.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_UNIT_CANNOT_HAVE_CHILDREN AdGroupCriterionErrorReason = "PRODUCT_PARTITION_UNIT_CANNOT_HAVE_CHILDREN"

	//
	// Subdivided product partitions must have an "others" case.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_SUBDIVISION_REQUIRES_OTHERS_CASE AdGroupCriterionErrorReason = "PRODUCT_PARTITION_SUBDIVISION_REQUIRES_OTHERS_CASE"

	//
	// Dimension type of product partition must be the same as that of its siblings.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_REQUIRES_SAME_DIMENSION_TYPE_AS_SIBLINGS AdGroupCriterionErrorReason = "PRODUCT_PARTITION_REQUIRES_SAME_DIMENSION_TYPE_AS_SIBLINGS"

	//
	// Product partition cannot be added to the ad group because it already exists.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_ALREADY_EXISTS AdGroupCriterionErrorReason = "PRODUCT_PARTITION_ALREADY_EXISTS"

	//
	// Product partition referenced in the operation was not found in the ad group.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_DOES_NOT_EXIST AdGroupCriterionErrorReason = "PRODUCT_PARTITION_DOES_NOT_EXIST"

	//
	// Recursive removal failed because product partition subdivision is being created or modified
	// in this request.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_CANNOT_BE_REMOVED AdGroupCriterionErrorReason = "PRODUCT_PARTITION_CANNOT_BE_REMOVED"

	//
	// Product partition type is not allowed for specified AdGroupCriterion type.
	//
	AdGroupCriterionErrorReasonINVALID_PRODUCT_PARTITION_TYPE AdGroupCriterionErrorReason = "INVALID_PRODUCT_PARTITION_TYPE"

	//
	// Product partition in an ADD operation specifies a non temporary CriterionId.
	//
	AdGroupCriterionErrorReasonPRODUCT_PARTITION_ADD_MAY_ONLY_USE_TEMP_ID AdGroupCriterionErrorReason = "PRODUCT_PARTITION_ADD_MAY_ONLY_USE_TEMP_ID"

	//
	// Partial failure is not supported for shopping campaign mutate operations.
	//
	AdGroupCriterionErrorReasonCAMPAIGN_TYPE_NOT_COMPATIBLE_WITH_PARTIAL_FAILURE AdGroupCriterionErrorReason = "CAMPAIGN_TYPE_NOT_COMPATIBLE_WITH_PARTIAL_FAILURE"

	//
	// Operations in the mutate request changes too many shopping ad groups. Please split
	// requests for multiple shopping ad groups across multiple requests.
	//
	AdGroupCriterionErrorReasonOPERATIONS_FOR_TOO_MANY_SHOPPING_ADGROUPS AdGroupCriterionErrorReason = "OPERATIONS_FOR_TOO_MANY_SHOPPING_ADGROUPS"

	//
	// Not allowed to modify url fields of an ad group criterion if there are duplicate elements
	// for that ad group criterion in the request.
	//
	AdGroupCriterionErrorReasonCANNOT_MODIFY_URL_FIELDS_WITH_DUPLICATE_ELEMENTS AdGroupCriterionErrorReason = "CANNOT_MODIFY_URL_FIELDS_WITH_DUPLICATE_ELEMENTS"

	//
	// Cannot set url fields without also setting final urls.
	//
	AdGroupCriterionErrorReasonCANNOT_SET_WITHOUT_FINAL_URLS AdGroupCriterionErrorReason = "CANNOT_SET_WITHOUT_FINAL_URLS"

	//
	// Cannot clear final urls if final mobile urls exist.
	//
	AdGroupCriterionErrorReasonCANNOT_CLEAR_FINAL_URLS_IF_FINAL_MOBILE_URLS_EXIST AdGroupCriterionErrorReason = "CANNOT_CLEAR_FINAL_URLS_IF_FINAL_MOBILE_URLS_EXIST"

	//
	// Cannot clear final urls if final app urls exist.
	//
	AdGroupCriterionErrorReasonCANNOT_CLEAR_FINAL_URLS_IF_FINAL_APP_URLS_EXIST AdGroupCriterionErrorReason = "CANNOT_CLEAR_FINAL_URLS_IF_FINAL_APP_URLS_EXIST"

	//
	// Cannot clear final urls if tracking url template exists.
	//
	AdGroupCriterionErrorReasonCANNOT_CLEAR_FINAL_URLS_IF_TRACKING_URL_TEMPLATE_EXISTS AdGroupCriterionErrorReason = "CANNOT_CLEAR_FINAL_URLS_IF_TRACKING_URL_TEMPLATE_EXISTS"

	//
	// Cannot clear final urls if url custom parameters exist.
	//
	AdGroupCriterionErrorReasonCANNOT_CLEAR_FINAL_URLS_IF_URL_CUSTOM_PARAMETERS_EXIST AdGroupCriterionErrorReason = "CANNOT_CLEAR_FINAL_URLS_IF_URL_CUSTOM_PARAMETERS_EXIST"

	//
	// Cannot set both destination url and final urls.
	//
	AdGroupCriterionErrorReasonCANNOT_SET_BOTH_DESTINATION_URL_AND_FINAL_URLS AdGroupCriterionErrorReason = "CANNOT_SET_BOTH_DESTINATION_URL_AND_FINAL_URLS"

	//
	// Cannot set both destination url and tracking url template.
	//
	AdGroupCriterionErrorReasonCANNOT_SET_BOTH_DESTINATION_URL_AND_TRACKING_URL_TEMPLATE AdGroupCriterionErrorReason = "CANNOT_SET_BOTH_DESTINATION_URL_AND_TRACKING_URL_TEMPLATE"

	//
	// Final urls are not supported for this criterion type.
	//
	AdGroupCriterionErrorReasonFINAL_URLS_NOT_SUPPORTED_FOR_CRITERION_TYPE AdGroupCriterionErrorReason = "FINAL_URLS_NOT_SUPPORTED_FOR_CRITERION_TYPE"

	//
	// Final mobile urls are not supported for this criterion type.
	//
	AdGroupCriterionErrorReasonFINAL_MOBILE_URLS_NOT_SUPPORTED_FOR_CRITERION_TYPE AdGroupCriterionErrorReason = "FINAL_MOBILE_URLS_NOT_SUPPORTED_FOR_CRITERION_TYPE"

	AdGroupCriterionErrorReasonUNKNOWN AdGroupCriterionErrorReason = "UNKNOWN"
)

//
// Error reasons.
//
type AdGroupFeedErrorReason string

const (

	//
	// An active feed already exists for this adgroup and place holder type.
	//
	AdGroupFeedErrorReasonFEED_ALREADY_EXISTS_FOR_PLACEHOLDER_TYPE AdGroupFeedErrorReason = "FEED_ALREADY_EXISTS_FOR_PLACEHOLDER_TYPE"

	//
	// The specified id does not exist.
	//
	AdGroupFeedErrorReasonINVALID_ID AdGroupFeedErrorReason = "INVALID_ID"

	//
	// The specified feed is deleted.
	//
	AdGroupFeedErrorReasonCANNOT_ADD_FOR_DELETED_FEED AdGroupFeedErrorReason = "CANNOT_ADD_FOR_DELETED_FEED"

	//
	// The AdGroupFeed already exists. SET should be used to modify the existing AdGroupFeed.
	//
	AdGroupFeedErrorReasonCANNOT_ADD_ALREADY_EXISTING_ADGROUP_FEED AdGroupFeedErrorReason = "CANNOT_ADD_ALREADY_EXISTING_ADGROUP_FEED"

	//
	// Cannot operate on removed adgroup feed.
	//
	AdGroupFeedErrorReasonCANNOT_OPERATE_ON_REMOVED_ADGROUP_FEED AdGroupFeedErrorReason = "CANNOT_OPERATE_ON_REMOVED_ADGROUP_FEED"

	//
	// Invalid placeholder type ids.
	//
	AdGroupFeedErrorReasonINVALID_PLACEHOLDER_TYPES AdGroupFeedErrorReason = "INVALID_PLACEHOLDER_TYPES"

	//
	// Feed mapping for this placeholder type does not exist.
	//
	AdGroupFeedErrorReasonMISSING_FEEDMAPPING_FOR_PLACEHOLDER_TYPE AdGroupFeedErrorReason = "MISSING_FEEDMAPPING_FOR_PLACEHOLDER_TYPE"

	//
	// Location AdGroupFeeds cannot be created unless there is a location CustomerFeed
	// for the specified feed.
	//
	AdGroupFeedErrorReasonNO_EXISTING_LOCATION_CUSTOMER_FEED AdGroupFeedErrorReason = "NO_EXISTING_LOCATION_CUSTOMER_FEED"

	AdGroupFeedErrorReasonUNKNOWN AdGroupFeedErrorReason = "UNKNOWN"
)

//
// The reasons for the adgroup service error.
//
type AdGroupServiceErrorReason string

const (

	//
	// AdGroup with the same name already exists for the campaign.
	//
	AdGroupServiceErrorReasonDUPLICATE_ADGROUP_NAME AdGroupServiceErrorReason = "DUPLICATE_ADGROUP_NAME"

	//
	// AdGroup name is not valid.
	//
	AdGroupServiceErrorReasonINVALID_ADGROUP_NAME AdGroupServiceErrorReason = "INVALID_ADGROUP_NAME"

	//
	// Cannot remove an adgroup, adgroup status can be marked removed
	// using set operator.
	//
	AdGroupServiceErrorReasonUSE_SET_OPERATOR_AND_MARK_STATUS_TO_REMOVED AdGroupServiceErrorReason = "USE_SET_OPERATOR_AND_MARK_STATUS_TO_REMOVED"

	//
	// Advertiser is not allowed to target sites or set site bids that are
	// not on the Google Search Network.
	//
	AdGroupServiceErrorReasonADVERTISER_NOT_ON_CONTENT_NETWORK AdGroupServiceErrorReason = "ADVERTISER_NOT_ON_CONTENT_NETWORK"

	//
	// Bid amount is too big.
	//
	AdGroupServiceErrorReasonBID_TOO_BIG AdGroupServiceErrorReason = "BID_TOO_BIG"

	//
	// AdGroup bid does not match the campaign's bidding strategy.
	//
	AdGroupServiceErrorReasonBID_TYPE_AND_BIDDING_STRATEGY_MISMATCH AdGroupServiceErrorReason = "BID_TYPE_AND_BIDDING_STRATEGY_MISMATCH"

	//
	// AdGroup name is required for Add.
	//
	AdGroupServiceErrorReasonMISSING_ADGROUP_NAME AdGroupServiceErrorReason = "MISSING_ADGROUP_NAME"

	//
	// No link found between the ad group and the label.
	//
	AdGroupServiceErrorReasonADGROUP_LABEL_DOES_NOT_EXIST AdGroupServiceErrorReason = "ADGROUP_LABEL_DOES_NOT_EXIST"

	//
	// The label has already been attached to the ad group.
	//
	AdGroupServiceErrorReasonADGROUP_LABEL_ALREADY_EXISTS AdGroupServiceErrorReason = "ADGROUP_LABEL_ALREADY_EXISTS"

	//
	// The CriterionTypeGroup is not supported for the content bid dimension.
	//
	AdGroupServiceErrorReasonINVALID_CONTENT_BID_CRITERION_TYPE_GROUP AdGroupServiceErrorReason = "INVALID_CONTENT_BID_CRITERION_TYPE_GROUP"

	//
	// The ad group type is not compatible with the campaign channel type.
	//
	AdGroupServiceErrorReasonAD_GROUP_TYPE_NOT_VALID_FOR_ADVERTISING_CHANNEL_TYPE AdGroupServiceErrorReason = "AD_GROUP_TYPE_NOT_VALID_FOR_ADVERTISING_CHANNEL_TYPE"

	//
	// The ad group type is not supported in the country of sale of the campaign.
	//
	AdGroupServiceErrorReasonADGROUP_TYPE_NOT_SUPPORTED_FOR_CAMPAIGN_SALES_COUNTRY AdGroupServiceErrorReason = "ADGROUP_TYPE_NOT_SUPPORTED_FOR_CAMPAIGN_SALES_COUNTRY"

	//
	// Ad groups of AdGroupType.SEARCH_DYNAMIC_ADS can only be added to campaigns that have
	// DynamicSearchAdsSetting attached.
	//
	AdGroupServiceErrorReasonCANNOT_ADD_ADGROUP_OF_TYPE_DSA_TO_CAMPAIGN_WITHOUT_DSA_SETTING AdGroupServiceErrorReason = "CANNOT_ADD_ADGROUP_OF_TYPE_DSA_TO_CAMPAIGN_WITHOUT_DSA_SETTING"
)

//
// Reasons for error.
//
type AdSharingErrorReason string

const (

	//
	// Error resulting in attempting to add an Ad to an AdGroup that already contains the Ad.
	//
	AdSharingErrorReasonAD_GROUP_ALREADY_CONTAINS_AD AdSharingErrorReason = "AD_GROUP_ALREADY_CONTAINS_AD"

	//
	// Ad is not compatible with the AdGroup it is being shared with. For more details, look
	// at {@link #sharedAdError}.
	//
	AdSharingErrorReasonINCOMPATIBLE_AD_UNDER_AD_GROUP AdSharingErrorReason = "INCOMPATIBLE_AD_UNDER_AD_GROUP"

	//
	// Cannot add AdGroupAd on inactive Ad.
	//
	AdSharingErrorReasonCANNOT_SHARE_INACTIVE_AD AdSharingErrorReason = "CANNOT_SHARE_INACTIVE_AD"
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

type BiddingErrorsReason string

const (

	//
	// Cannot transition to new bidding strategy.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_TRANSITION_NOT_ALLOWED BiddingErrorsReason = "BIDDING_STRATEGY_TRANSITION_NOT_ALLOWED"

	//
	// Bidding strategy cannot be overridden by current ad group bidding strategy.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_NOT_COMPATIBLE_WITH_ADGROUP_OVERRIDES BiddingErrorsReason = "BIDDING_STRATEGY_NOT_COMPATIBLE_WITH_ADGROUP_OVERRIDES"

	//
	// Bidding strategy cannot be overridden by current criteria bidding strategy.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_NOT_COMPATIBLE_WITH_ADGROUP_CRITERIA_OVERRIDES BiddingErrorsReason = "BIDDING_STRATEGY_NOT_COMPATIBLE_WITH_ADGROUP_CRITERIA_OVERRIDES"

	//
	// Cannot override campaign bidding strategy.
	//
	BiddingErrorsReasonCAMPAIGN_BIDDING_STRATEGY_CANNOT_BE_OVERRIDDEN BiddingErrorsReason = "CAMPAIGN_BIDDING_STRATEGY_CANNOT_BE_OVERRIDDEN"

	//
	// Cannot override ad group bidding strategy.
	//
	BiddingErrorsReasonADGROUP_BIDDING_STRATEGY_CANNOT_BE_OVERRIDDEN BiddingErrorsReason = "ADGROUP_BIDDING_STRATEGY_CANNOT_BE_OVERRIDDEN"

	//
	// Cannot attach bidding strategy to campaign.
	//
	BiddingErrorsReasonCANNOT_ATTACH_BIDDING_STRATEGY_TO_CAMPAIGN BiddingErrorsReason = "CANNOT_ATTACH_BIDDING_STRATEGY_TO_CAMPAIGN"

	//
	// Cannot attach bidding strategy to ad group.
	//
	BiddingErrorsReasonCANNOT_ATTACH_BIDDING_STRATEGY_TO_ADGROUP BiddingErrorsReason = "CANNOT_ATTACH_BIDDING_STRATEGY_TO_ADGROUP"

	//
	// Cannot attach bidding strategy to criteria.
	//
	BiddingErrorsReasonCANNOT_ATTACH_BIDDING_STRATEGY_TO_ADGROUP_CRITERIA BiddingErrorsReason = "CANNOT_ATTACH_BIDDING_STRATEGY_TO_ADGROUP_CRITERIA"

	//
	// Bidding strategy is not supported or cannot be used as anonymous.
	//
	BiddingErrorsReasonINVALID_ANONYMOUS_BIDDING_STRATEGY_TYPE BiddingErrorsReason = "INVALID_ANONYMOUS_BIDDING_STRATEGY_TYPE"

	//
	// No bids may be set. The bid list must be null or empty.
	//
	BiddingErrorsReasonBIDS_NOT_ALLLOWED BiddingErrorsReason = "BIDS_NOT_ALLLOWED"

	//
	// The bid list contains two or more bids of the same type.
	//
	BiddingErrorsReasonDUPLICATE_BIDS BiddingErrorsReason = "DUPLICATE_BIDS"

	//
	// The bidding scheme does not match the bidding strategy type.
	//
	BiddingErrorsReasonINVALID_BIDDING_SCHEME BiddingErrorsReason = "INVALID_BIDDING_SCHEME"

	//
	// The type does not match the named strategy's type.
	//
	BiddingErrorsReasonINVALID_BIDDING_STRATEGY_TYPE BiddingErrorsReason = "INVALID_BIDDING_STRATEGY_TYPE"

	//
	// The bidding strategy type is missing.
	//
	BiddingErrorsReasonMISSING_BIDDING_STRATEGY_TYPE BiddingErrorsReason = "MISSING_BIDDING_STRATEGY_TYPE"

	//
	// The bid list contains a null entry.
	//
	BiddingErrorsReasonNULL_BID BiddingErrorsReason = "NULL_BID"

	//
	// The bid is invalid.
	//
	BiddingErrorsReasonINVALID_BID BiddingErrorsReason = "INVALID_BID"

	//
	// Bidding strategy is not available for the account type.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_NOT_AVAILABLE_FOR_ACCOUNT_TYPE BiddingErrorsReason = "BIDDING_STRATEGY_NOT_AVAILABLE_FOR_ACCOUNT_TYPE"

	//
	// Conversion tracking is not enabled for the campaign for VBB transition.
	//
	BiddingErrorsReasonCONVERSION_TRACKING_NOT_ENABLED BiddingErrorsReason = "CONVERSION_TRACKING_NOT_ENABLED"

	//
	// Not enough conversions tracked for VBB transitions.
	//
	BiddingErrorsReasonNOT_ENOUGH_CONVERSIONS BiddingErrorsReason = "NOT_ENOUGH_CONVERSIONS"

	//
	// Campaign can not be created with given bidding strategy. It can be transitioned to the
	// strategy, once eligible.
	//
	BiddingErrorsReasonCANNOT_CREATE_CAMPAIGN_WITH_BIDDING_STRATEGY BiddingErrorsReason = "CANNOT_CREATE_CAMPAIGN_WITH_BIDDING_STRATEGY"

	//
	// Cannot target content network only as ad group uses Page One Promoted bidding strategy.
	//
	BiddingErrorsReasonCANNOT_TARGET_CONTENT_NETWORK_ONLY_WITH_AD_GROUP_LEVEL_POP_BIDDING_STRATEGY BiddingErrorsReason = "CANNOT_TARGET_CONTENT_NETWORK_ONLY_WITH_AD_GROUP_LEVEL_POP_BIDDING_STRATEGY"

	//
	// Cannot target content network only as campaign uses Page One Promoted bidding strategy.
	//
	BiddingErrorsReasonCANNOT_TARGET_CONTENT_NETWORK_ONLY_WITH_CAMPAIGN_LEVEL_POP_BIDDING_STRATEGY BiddingErrorsReason = "CANNOT_TARGET_CONTENT_NETWORK_ONLY_WITH_CAMPAIGN_LEVEL_POP_BIDDING_STRATEGY"

	//
	// Budget Optimizer and Target Spend bidding strategies are not supported for campaigns with
	// AdSchedule targeting.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_NOT_SUPPORTED_WITH_AD_SCHEDULE BiddingErrorsReason = "BIDDING_STRATEGY_NOT_SUPPORTED_WITH_AD_SCHEDULE"

	//
	// Pay per conversion is not available to all the customer, only few whitelisted customers
	// can use this.
	//
	BiddingErrorsReasonPAY_PER_CONVERSION_NOT_AVAILABLE_FOR_CUSTOMER BiddingErrorsReason = "PAY_PER_CONVERSION_NOT_AVAILABLE_FOR_CUSTOMER"

	//
	// Pay per conversion is not allowed with Target CPA.
	//
	BiddingErrorsReasonPAY_PER_CONVERSION_NOT_ALLOWED_WITH_TARGET_CPA BiddingErrorsReason = "PAY_PER_CONVERSION_NOT_ALLOWED_WITH_TARGET_CPA"

	//
	// Cannot set bidding strategy to Manual CPM for search network only campaigns.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_NOT_ALLOWED_FOR_SEARCH_ONLY_CAMPAIGNS BiddingErrorsReason = "BIDDING_STRATEGY_NOT_ALLOWED_FOR_SEARCH_ONLY_CAMPAIGNS"

	//
	// The bidding strategy is not supported for use in drafts or experiments.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_NOT_SUPPORTED_IN_DRAFTS_OR_EXPERIMENTS BiddingErrorsReason = "BIDDING_STRATEGY_NOT_SUPPORTED_IN_DRAFTS_OR_EXPERIMENTS"

	//
	// Bidding strategy type does not support product type ad group criterion.
	//
	BiddingErrorsReasonBIDDING_STRATEGY_TYPE_DOES_NOT_SUPPORT_PRODUCT_TYPE_ADGROUP_CRITERION BiddingErrorsReason = "BIDDING_STRATEGY_TYPE_DOES_NOT_SUPPORT_PRODUCT_TYPE_ADGROUP_CRITERION"

	//
	// Bid amount is too small.
	//
	BiddingErrorsReasonBID_TOO_SMALL BiddingErrorsReason = "BID_TOO_SMALL"

	//
	// Bid amount is too big.
	//
	BiddingErrorsReasonBID_TOO_BIG BiddingErrorsReason = "BID_TOO_BIG"

	//
	// Bid has too many fractional digit precision.
	//
	BiddingErrorsReasonBID_TOO_MANY_FRACTIONAL_DIGITS BiddingErrorsReason = "BID_TOO_MANY_FRACTIONAL_DIGITS"

	//
	// EnhancedCpcEnabled cannot be set on <em>portfolio</em> Target Spend strategies.
	//
	BiddingErrorsReasonENHANCED_CPC_ENABLED_NOT_SUPPORTED_ON_PORTFOLIO_TARGET_SPEND_STRATEGY BiddingErrorsReason = "ENHANCED_CPC_ENABLED_NOT_SUPPORTED_ON_PORTFOLIO_TARGET_SPEND_STRATEGY"

	BiddingErrorsReasonUNKNOWN BiddingErrorsReason = "UNKNOWN"
)

//
// The reasons for the error.
//
type CampaignBidModifierErrorReason string

const (
	CampaignBidModifierErrorReasonCAMPAIGN_BID_MODIFIER_ERROR CampaignBidModifierErrorReason = "CAMPAIGN_BID_MODIFIER_ERROR"
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
// The reasons for the target error.
//
type CampaignErrorReason string

const (

	//
	// A complete campaign cannot go back to being incomplete
	//
	CampaignErrorReasonCANNOT_GO_BACK_TO_INCOMPLETE CampaignErrorReason = "CANNOT_GO_BACK_TO_INCOMPLETE"

	//
	// Cannot target content network.
	//
	CampaignErrorReasonCANNOT_TARGET_CONTENT_NETWORK CampaignErrorReason = "CANNOT_TARGET_CONTENT_NETWORK"

	//
	// Cannot target search network.
	//
	CampaignErrorReasonCANNOT_TARGET_SEARCH_NETWORK CampaignErrorReason = "CANNOT_TARGET_SEARCH_NETWORK"

	//
	// Cannot cover search network without google search network.
	//
	CampaignErrorReasonCANNOT_TARGET_SEARCH_NETWORK_WITHOUT_GOOGLE_SEARCH CampaignErrorReason = "CANNOT_TARGET_SEARCH_NETWORK_WITHOUT_GOOGLE_SEARCH"

	//
	// Cannot target Google Search network for a CPM campaign.
	//
	CampaignErrorReasonCANNOT_TARGET_GOOGLE_SEARCH_FOR_CPM_CAMPAIGN CampaignErrorReason = "CANNOT_TARGET_GOOGLE_SEARCH_FOR_CPM_CAMPAIGN"

	//
	// Must target at least one network.
	//
	CampaignErrorReasonCAMPAIGN_MUST_TARGET_AT_LEAST_ONE_NETWORK CampaignErrorReason = "CAMPAIGN_MUST_TARGET_AT_LEAST_ONE_NETWORK"

	//
	// Only some Google partners are allowed to target partner search network.
	//
	CampaignErrorReasonCANNOT_TARGET_PARTNER_SEARCH_NETWORK CampaignErrorReason = "CANNOT_TARGET_PARTNER_SEARCH_NETWORK"

	//
	// Cannot target content network only as campaign has criteria-level bidding strategy.
	//
	CampaignErrorReasonCANNOT_TARGET_CONTENT_NETWORK_ONLY_WITH_CRITERIA_LEVEL_BIDDING_STRATEGY CampaignErrorReason = "CANNOT_TARGET_CONTENT_NETWORK_ONLY_WITH_CRITERIA_LEVEL_BIDDING_STRATEGY"

	//
	// Cannot modify the start or end date such that the campaign duration would not contain the
	// durations of all runnable trials.
	//
	CampaignErrorReasonCAMPAIGN_DURATION_MUST_CONTAIN_ALL_RUNNABLE_TRIALS CampaignErrorReason = "CAMPAIGN_DURATION_MUST_CONTAIN_ALL_RUNNABLE_TRIALS"

	//
	// Cannot modify dates, budget or campaign name of a trial campaign.
	//
	CampaignErrorReasonCANNOT_MODIFY_FOR_TRIAL_CAMPAIGN CampaignErrorReason = "CANNOT_MODIFY_FOR_TRIAL_CAMPAIGN"

	//
	// Trying to modify the name of an active or paused campaign, where the name is already
	// assigned to another active or paused campaign.
	//
	CampaignErrorReasonDUPLICATE_CAMPAIGN_NAME CampaignErrorReason = "DUPLICATE_CAMPAIGN_NAME"

	//
	// Two fields are in conflicting modes.
	//
	CampaignErrorReasonINCOMPATIBLE_CAMPAIGN_FIELD CampaignErrorReason = "INCOMPATIBLE_CAMPAIGN_FIELD"

	//
	// Campaign name cannot be used.
	//
	CampaignErrorReasonINVALID_CAMPAIGN_NAME CampaignErrorReason = "INVALID_CAMPAIGN_NAME"

	//
	// Given status is invalid.
	//
	CampaignErrorReasonINVALID_AD_SERVING_OPTIMIZATION_STATUS CampaignErrorReason = "INVALID_AD_SERVING_OPTIMIZATION_STATUS"

	//
	// Error in the campaign level tracking url.
	//
	CampaignErrorReasonINVALID_TRACKING_URL CampaignErrorReason = "INVALID_TRACKING_URL"

	//
	// Cannot set both tracking url template and tracking setting. An user has to clear legacy
	// tracking setting in order to add tracking url template.
	//
	CampaignErrorReasonCANNOT_SET_BOTH_TRACKING_URL_TEMPLATE_AND_TRACKING_SETTING CampaignErrorReason = "CANNOT_SET_BOTH_TRACKING_URL_TEMPLATE_AND_TRACKING_SETTING"

	//
	// The maximum number of impressions for Frequency Cap should be an integer greater than 0.
	//
	CampaignErrorReasonMAX_IMPRESSIONS_NOT_IN_RANGE CampaignErrorReason = "MAX_IMPRESSIONS_NOT_IN_RANGE"

	//
	// Only the Day, Week and Month time units are supported.
	//
	CampaignErrorReasonTIME_UNIT_NOT_SUPPORTED CampaignErrorReason = "TIME_UNIT_NOT_SUPPORTED"

	//
	// Operation not allowed on a campaign whose serving status has ended
	//
	CampaignErrorReasonINVALID_OPERATION_IF_SERVING_STATUS_HAS_ENDED CampaignErrorReason = "INVALID_OPERATION_IF_SERVING_STATUS_HAS_ENDED"

	//
	// This budget is exclusively linked to a Campaign that is using @link{Experiment}s
	// so it cannot be shared.
	//
	CampaignErrorReasonBUDGET_CANNOT_BE_SHARED CampaignErrorReason = "BUDGET_CANNOT_BE_SHARED"

	//
	// Campaigns using @link{Experiment}s cannot use a shared budget.
	//
	CampaignErrorReasonCAMPAIGN_CANNOT_USE_SHARED_BUDGET CampaignErrorReason = "CAMPAIGN_CANNOT_USE_SHARED_BUDGET"

	//
	// A different budget cannot be assigned to a campaign when there are running or scheduled
	// trials.
	//
	CampaignErrorReasonCANNOT_CHANGE_BUDGET_ON_CAMPAIGN_WITH_TRIALS CampaignErrorReason = "CANNOT_CHANGE_BUDGET_ON_CAMPAIGN_WITH_TRIALS"

	//
	// No link found between the campaign and the label.
	//
	CampaignErrorReasonCAMPAIGN_LABEL_DOES_NOT_EXIST CampaignErrorReason = "CAMPAIGN_LABEL_DOES_NOT_EXIST"

	//
	// The label has already been attached to the campaign.
	//
	CampaignErrorReasonCAMPAIGN_LABEL_ALREADY_EXISTS CampaignErrorReason = "CAMPAIGN_LABEL_ALREADY_EXISTS"

	//
	// A ShoppingSetting was not found when creating a shopping campaign.
	//
	CampaignErrorReasonMISSING_SHOPPING_SETTING CampaignErrorReason = "MISSING_SHOPPING_SETTING"

	//
	// The country in shopping setting is not an allowed country.
	//
	CampaignErrorReasonINVALID_SHOPPING_SALES_COUNTRY CampaignErrorReason = "INVALID_SHOPPING_SALES_COUNTRY"

	//
	// Shopping merchant is not enabled for Purchases on Google.
	//
	CampaignErrorReasonSHOPPING_MERCHANT_NOT_ALLOWED_FOR_PURCHASES_ON_GOOGLE CampaignErrorReason = "SHOPPING_MERCHANT_NOT_ALLOWED_FOR_PURCHASES_ON_GOOGLE"

	//
	// Purchases on Google not enabled for the shopping campaign's sales country.
	//
	CampaignErrorReasonPURCHASES_ON_GOOGLE_NOT_SUPPORTED_FOR_SHOPPING_SALES_COUNTRY CampaignErrorReason = "PURCHASES_ON_GOOGLE_NOT_SUPPORTED_FOR_SHOPPING_SALES_COUNTRY"

	//
	// A Campaign with channel sub type UNIVERSAL_APP_CAMPAIGN must have a
	// UniversalAppCampaignSetting specified.
	//
	CampaignErrorReasonMISSING_UNIVERSAL_APP_CAMPAIGN_SETTING CampaignErrorReason = "MISSING_UNIVERSAL_APP_CAMPAIGN_SETTING"

	//
	// The requested channel type is not available according to the customer's account setting.
	//
	CampaignErrorReasonADVERTISING_CHANNEL_TYPE_NOT_AVAILABLE_FOR_ACCOUNT_TYPE CampaignErrorReason = "ADVERTISING_CHANNEL_TYPE_NOT_AVAILABLE_FOR_ACCOUNT_TYPE"

	//
	// The AdvertisingChannelSubType is not a valid subtype of the primary channel type.
	//
	CampaignErrorReasonINVALID_ADVERTISING_CHANNEL_SUB_TYPE CampaignErrorReason = "INVALID_ADVERTISING_CHANNEL_SUB_TYPE"

	//
	// At least one conversion must be selected.
	//
	CampaignErrorReasonAT_LEAST_ONE_CONVERSION_MUST_BE_SELECTED CampaignErrorReason = "AT_LEAST_ONE_CONVERSION_MUST_BE_SELECTED"

	//
	// Setting ad rotation mode for a campaign is not allowed.
	// Ad rotation mode at campaign is deprecated.
	//
	CampaignErrorReasonCANNOT_SET_AD_ROTATION_MODE CampaignErrorReason = "CANNOT_SET_AD_ROTATION_MODE"

	//
	// Default error
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	CampaignErrorReasonUNKNOWN CampaignErrorReason = "UNKNOWN"
)

//
// Error reasons.
//
type CampaignFeedErrorReason string

const (

	//
	// An active feed already exists for this campaign and place holder type.
	//
	CampaignFeedErrorReasonFEED_ALREADY_EXISTS_FOR_PLACEHOLDER_TYPE CampaignFeedErrorReason = "FEED_ALREADY_EXISTS_FOR_PLACEHOLDER_TYPE"

	//
	// The specified id does not exist.
	//
	CampaignFeedErrorReasonINVALID_ID CampaignFeedErrorReason = "INVALID_ID"

	//
	// The specified feed is deleted.
	//
	CampaignFeedErrorReasonCANNOT_ADD_FOR_DELETED_FEED CampaignFeedErrorReason = "CANNOT_ADD_FOR_DELETED_FEED"

	//
	// The CampaignFeed already exists. SET should be used to modify the existing CampaignFeed.
	//
	CampaignFeedErrorReasonCANNOT_ADD_ALREADY_EXISTING_CAMPAIGN_FEED CampaignFeedErrorReason = "CANNOT_ADD_ALREADY_EXISTING_CAMPAIGN_FEED"

	//
	// Cannot operate on deleted campaign feed.
	//
	CampaignFeedErrorReasonCANNOT_OPERATE_ON_REMOVED_CAMPAIGN_FEED CampaignFeedErrorReason = "CANNOT_OPERATE_ON_REMOVED_CAMPAIGN_FEED"

	//
	// Invalid placeholder type ids.
	//
	CampaignFeedErrorReasonINVALID_PLACEHOLDER_TYPES CampaignFeedErrorReason = "INVALID_PLACEHOLDER_TYPES"

	//
	// Feed mapping for this placeholder type does not exist.
	//
	CampaignFeedErrorReasonMISSING_FEEDMAPPING_FOR_PLACEHOLDER_TYPE CampaignFeedErrorReason = "MISSING_FEEDMAPPING_FOR_PLACEHOLDER_TYPE"

	//
	// Location CampaignFeeds cannot be created unless there is a location CustomerFeed
	// for the specified feed.
	//
	CampaignFeedErrorReasonNO_EXISTING_LOCATION_CUSTOMER_FEED CampaignFeedErrorReason = "NO_EXISTING_LOCATION_CUSTOMER_FEED"

	CampaignFeedErrorReasonUNKNOWN CampaignFeedErrorReason = "UNKNOWN"
)

type CampaignPreferenceErrorReason string

const (

	//
	// A campaign cannot have two preferences with the same preference key.
	//
	CampaignPreferenceErrorReasonPREFERENCE_ALREADY_EXISTS CampaignPreferenceErrorReason = "PREFERENCE_ALREADY_EXISTS"

	//
	// No preference matched the given preference key.
	//
	CampaignPreferenceErrorReasonPREFERENCE_NOT_FOUND CampaignPreferenceErrorReason = "PREFERENCE_NOT_FOUND"

	CampaignPreferenceErrorReasonUNKNOWN CampaignPreferenceErrorReason = "UNKNOWN"
)

//
// Error reasons
//
type CampaignSharedSetErrorReason string

const (
	CampaignSharedSetErrorReasonCAMPAIGN_SHARED_SET_DOES_NOT_EXIST CampaignSharedSetErrorReason = "CAMPAIGN_SHARED_SET_DOES_NOT_EXIST"

	CampaignSharedSetErrorReasonUNKNOWN CampaignSharedSetErrorReason = "UNKNOWN"
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
// The reasons for the date range error.
//
type DateRangeErrorReason string

const (
	DateRangeErrorReasonDATE_RANGE_ERROR DateRangeErrorReason = "DATE_RANGE_ERROR"

	//
	// Invalid date.
	//
	DateRangeErrorReasonINVALID_DATE DateRangeErrorReason = "INVALID_DATE"

	//
	// The start date was after the end date.
	//
	DateRangeErrorReasonSTART_DATE_AFTER_END_DATE DateRangeErrorReason = "START_DATE_AFTER_END_DATE"

	//
	// Cannot set date to past time
	//
	DateRangeErrorReasonCANNOT_SET_DATE_TO_PAST DateRangeErrorReason = "CANNOT_SET_DATE_TO_PAST"

	//
	// A date was used that is past the system "last" date.
	//
	DateRangeErrorReasonAFTER_MAXIMUM_ALLOWABLE_DATE DateRangeErrorReason = "AFTER_MAXIMUM_ALLOWABLE_DATE"

	//
	// Trying to change start date on a campaign that has started.
	//
	DateRangeErrorReasonCANNOT_MODIFY_START_DATE_IF_ALREADY_STARTED DateRangeErrorReason = "CANNOT_MODIFY_START_DATE_IF_ALREADY_STARTED"
)

//
// The reasons for the validation error.
//
type DistinctErrorReason string

const (
	DistinctErrorReasonDUPLICATE_ELEMENT DistinctErrorReason = "DUPLICATE_ELEMENT"

	DistinctErrorReasonDUPLICATE_TYPE DistinctErrorReason = "DUPLICATE_TYPE"
)

type DraftErrorReason string

const (

	//
	// The draft is archived and cannot be modified further.
	//
	DraftErrorReasonCANNOT_CHANGE_ARCHIVED_DRAFT DraftErrorReason = "CANNOT_CHANGE_ARCHIVED_DRAFT"

	//
	// The draft has been promoted and cannot be modified further.
	//
	DraftErrorReasonCANNOT_CHANGE_PROMOTED_DRAFT DraftErrorReason = "CANNOT_CHANGE_PROMOTED_DRAFT"

	//
	// The draft has failed to be promoted and cannot be modified further.
	//
	DraftErrorReasonCANNOT_CHANGE_PROMOTE_FAILED_DRAFT DraftErrorReason = "CANNOT_CHANGE_PROMOTE_FAILED_DRAFT"

	//
	// This customer is not allowed to create drafts.
	//
	DraftErrorReasonCUSTOMER_CANNOT_CREATE_DRAFT DraftErrorReason = "CUSTOMER_CANNOT_CREATE_DRAFT"

	//
	// This campaign is not allowed to create drafts.
	//
	DraftErrorReasonCAMPAIGN_CANNOT_CREATE_DRAFT DraftErrorReason = "CAMPAIGN_CANNOT_CREATE_DRAFT"

	//
	// A draft with this name already exists.
	//
	DraftErrorReasonDUPLICATE_DRAFT_NAME DraftErrorReason = "DUPLICATE_DRAFT_NAME"

	//
	// This modification cannot be made on a draft.
	//
	DraftErrorReasonINVALID_DRAFT_CHANGE DraftErrorReason = "INVALID_DRAFT_CHANGE"

	//
	// The draft cannot be transitioned to the specified status from the its current status.
	//
	DraftErrorReasonINVALID_STATUS_TRANSITION DraftErrorReason = "INVALID_STATUS_TRANSITION"

	//
	// The campaign has reached the maximum number of drafts that can be created for a campaign
	// throughout its lifetime. No additional drafts can be created for this campaign. Archived
	// drafts also count towards this limit.
	//
	DraftErrorReasonMAX_NUMBER_OF_DRAFTS_PER_CAMPAIGN_REACHED DraftErrorReason = "MAX_NUMBER_OF_DRAFTS_PER_CAMPAIGN_REACHED"

	DraftErrorReasonDRAFT_ERROR DraftErrorReason = "DRAFT_ERROR"
)

type EntityAccessDeniedReason string

const (

	//
	// User did not have read access.
	//
	EntityAccessDeniedReasonREAD_ACCESS_DENIED EntityAccessDeniedReason = "READ_ACCESS_DENIED"

	//
	// User did not have write access.
	//
	EntityAccessDeniedReasonWRITE_ACCESS_DENIED EntityAccessDeniedReason = "WRITE_ACCESS_DENIED"
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
// Error reasons.
//
type FeedErrorReason string

const (

	//
	// The names of the FeedAttributes must be unique.
	//
	FeedErrorReasonATTRIBUTE_NAMES_NOT_UNIQUE FeedErrorReason = "ATTRIBUTE_NAMES_NOT_UNIQUE"

	//
	// The attribute list must be an exact copy of the existing list if the attribute id's are
	// present.
	//
	FeedErrorReasonATTRIBUTES_DO_NOT_MATCH_EXISTING_ATTRIBUTES FeedErrorReason = "ATTRIBUTES_DO_NOT_MATCH_EXISTING_ATTRIBUTES"

	//
	// Origin can only be set during Feed creation.
	//
	FeedErrorReasonCANNOT_CHANGE_ORIGIN FeedErrorReason = "CANNOT_CHANGE_ORIGIN"

	//
	// Cannot specify USER origin for a system generated feed.
	//
	FeedErrorReasonCANNOT_SPECIFY_USER_ORIGIN_FOR_SYSTEM_FEED FeedErrorReason = "CANNOT_SPECIFY_USER_ORIGIN_FOR_SYSTEM_FEED"

	//
	// Cannot specify ADWORDS origin for a non-system generated feed.
	//
	FeedErrorReasonCANNOT_SPECIFY_ADWORDS_ORIGIN_FOR_NON_SYSTEM_FEED FeedErrorReason = "CANNOT_SPECIFY_ADWORDS_ORIGIN_FOR_NON_SYSTEM_FEED"

	//
	// Cannot specify feed attributes for system feed.
	//
	FeedErrorReasonCANNOT_SPECIFY_FEED_ATTRIBUTES_FOR_SYSTEM_FEED FeedErrorReason = "CANNOT_SPECIFY_FEED_ATTRIBUTES_FOR_SYSTEM_FEED"

	//
	// Cannot update FeedAttributes on feed with origin adwords.
	//
	FeedErrorReasonCANNOT_UPDATE_FEED_ATTRIBUTES_WITH_ORIGIN_ADWORDS FeedErrorReason = "CANNOT_UPDATE_FEED_ATTRIBUTES_WITH_ORIGIN_ADWORDS"

	//
	// The given id refers to a removed Feed. Removed Feeds are immutable.
	//
	FeedErrorReasonFEED_REMOVED FeedErrorReason = "FEED_REMOVED"

	//
	// The origin of the feed is not valid for the client.
	//
	FeedErrorReasonINVALID_ORIGIN_VALUE FeedErrorReason = "INVALID_ORIGIN_VALUE"

	//
	// A user can only create and modify feeds with user origin.
	//
	FeedErrorReasonFEED_ORIGIN_IS_NOT_USER FeedErrorReason = "FEED_ORIGIN_IS_NOT_USER"

	//
	// Invalid auth token for the given email
	//
	FeedErrorReasonINVALID_AUTH_TOKEN_FOR_EMAIL FeedErrorReason = "INVALID_AUTH_TOKEN_FOR_EMAIL"

	//
	// Invalid email specified
	//
	FeedErrorReasonINVALID_EMAIL FeedErrorReason = "INVALID_EMAIL"

	//
	// Feed name matches that of another active Feed.
	//
	FeedErrorReasonDUPLICATE_FEED_NAME FeedErrorReason = "DUPLICATE_FEED_NAME"

	//
	// Name of feed is not allowed.
	//
	FeedErrorReasonINVALID_FEED_NAME FeedErrorReason = "INVALID_FEED_NAME"

	//
	// Missing OAuthInfo
	//
	FeedErrorReasonMISSING_OAUTH_INFO FeedErrorReason = "MISSING_OAUTH_INFO"

	//
	// New FeedAttributes must not effect the unique key.
	//
	FeedErrorReasonNEW_ATTRIBUTE_CANNOT_BE_PART_OF_UNIQUE_KEY FeedErrorReason = "NEW_ATTRIBUTE_CANNOT_BE_PART_OF_UNIQUE_KEY"

	//
	// Too many FeedAttributes for a Feed.
	//
	FeedErrorReasonTOO_MANY_FEED_ATTRIBUTES_FOR_FEED FeedErrorReason = "TOO_MANY_FEED_ATTRIBUTES_FOR_FEED"

	//
	// The business account is not valid.
	//
	FeedErrorReasonINVALID_BUSINESS_ACCOUNT FeedErrorReason = "INVALID_BUSINESS_ACCOUNT"

	//
	// Business account cannot access Google My Business account.
	//
	FeedErrorReasonBUSINESS_ACCOUNT_CANNOT_ACCESS_LOCATION_ACCOUNT FeedErrorReason = "BUSINESS_ACCOUNT_CANNOT_ACCESS_LOCATION_ACCOUNT"

	//
	// Invalid chain id provided for affiliate location feed.
	//
	FeedErrorReasonINVALID_AFFILIATE_CHAIN_ID FeedErrorReason = "INVALID_AFFILIATE_CHAIN_ID"

	//
	// Cannot change system feed generation data type
	//
	FeedErrorReasonCANNOT_CHANGE_SYSTEM_FEED_GENERATION_DATA_TYPE FeedErrorReason = "CANNOT_CHANGE_SYSTEM_FEED_GENERATION_DATA_TYPE"

	//
	// Unsupported relationship type
	//
	FeedErrorReasonUNSUPPORTED_AFFILIATE_LOCATION_RELATIONSHIP_TYPE FeedErrorReason = "UNSUPPORTED_AFFILIATE_LOCATION_RELATIONSHIP_TYPE"

	FeedErrorReasonUNKNOWN FeedErrorReason = "UNKNOWN"
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

type ImageErrorReason string

const (

	//
	// The image is not valid.
	//
	ImageErrorReasonINVALID_IMAGE ImageErrorReason = "INVALID_IMAGE"

	//
	// The image could not be stored.
	//
	ImageErrorReasonSTORAGE_ERROR ImageErrorReason = "STORAGE_ERROR"

	//
	// There was a problem with the request.
	//
	ImageErrorReasonBAD_REQUEST ImageErrorReason = "BAD_REQUEST"

	//
	// The image is not of legal dimensions.
	//
	ImageErrorReasonUNEXPECTED_SIZE ImageErrorReason = "UNEXPECTED_SIZE"

	//
	// Animated image are not permitted.
	//
	ImageErrorReasonANIMATED_NOT_ALLOWED ImageErrorReason = "ANIMATED_NOT_ALLOWED"

	//
	// Animation is too long.
	//
	ImageErrorReasonANIMATION_TOO_LONG ImageErrorReason = "ANIMATION_TOO_LONG"

	//
	// There was an error on the server.
	//
	ImageErrorReasonSERVER_ERROR ImageErrorReason = "SERVER_ERROR"

	//
	// Image cannot be in CMYK color format.
	//
	ImageErrorReasonCMYK_JPEG_NOT_ALLOWED ImageErrorReason = "CMYK_JPEG_NOT_ALLOWED"

	//
	// Flash images are not permitted.
	//
	ImageErrorReasonFLASH_NOT_ALLOWED ImageErrorReason = "FLASH_NOT_ALLOWED"

	//
	// Flash images must support clickTag.
	//
	ImageErrorReasonFLASH_WITHOUT_CLICKTAG ImageErrorReason = "FLASH_WITHOUT_CLICKTAG"

	//
	// A flash error has occurred after fixing the click tag.
	//
	ImageErrorReasonFLASH_ERROR_AFTER_FIXING_CLICK_TAG ImageErrorReason = "FLASH_ERROR_AFTER_FIXING_CLICK_TAG"

	//
	// Unacceptable visual effects.
	//
	ImageErrorReasonANIMATED_VISUAL_EFFECT ImageErrorReason = "ANIMATED_VISUAL_EFFECT"

	//
	// There was a problem with the flash image.
	//
	ImageErrorReasonFLASH_ERROR ImageErrorReason = "FLASH_ERROR"

	//
	// Incorrect image layout.
	//
	ImageErrorReasonLAYOUT_PROBLEM ImageErrorReason = "LAYOUT_PROBLEM"

	//
	// There was a problem reading the image file.
	//
	ImageErrorReasonPROBLEM_READING_IMAGE_FILE ImageErrorReason = "PROBLEM_READING_IMAGE_FILE"

	//
	// There was an error storing the image.
	//
	ImageErrorReasonERROR_STORING_IMAGE ImageErrorReason = "ERROR_STORING_IMAGE"

	//
	// The aspect ratio of the image is not allowed.
	//
	ImageErrorReasonASPECT_RATIO_NOT_ALLOWED ImageErrorReason = "ASPECT_RATIO_NOT_ALLOWED"

	//
	// Flash cannot have network objects.
	//
	ImageErrorReasonFLASH_HAS_NETWORK_OBJECTS ImageErrorReason = "FLASH_HAS_NETWORK_OBJECTS"

	//
	// Flash cannot have network methods.
	//
	ImageErrorReasonFLASH_HAS_NETWORK_METHODS ImageErrorReason = "FLASH_HAS_NETWORK_METHODS"

	//
	// Flash cannot have a Url.
	//
	ImageErrorReasonFLASH_HAS_URL ImageErrorReason = "FLASH_HAS_URL"

	//
	// Flash cannot use mouse tracking.
	//
	ImageErrorReasonFLASH_HAS_MOUSE_TRACKING ImageErrorReason = "FLASH_HAS_MOUSE_TRACKING"

	//
	// Flash cannot have a random number.
	//
	ImageErrorReasonFLASH_HAS_RANDOM_NUM ImageErrorReason = "FLASH_HAS_RANDOM_NUM"

	//
	// Ad click target cannot be '_self'.
	//
	ImageErrorReasonFLASH_SELF_TARGETS ImageErrorReason = "FLASH_SELF_TARGETS"

	//
	// GetUrl method should only use '_blank'.
	//
	ImageErrorReasonFLASH_BAD_GETURL_TARGET ImageErrorReason = "FLASH_BAD_GETURL_TARGET"

	//
	// Flash version is not supported.
	//
	ImageErrorReasonFLASH_VERSION_NOT_SUPPORTED ImageErrorReason = "FLASH_VERSION_NOT_SUPPORTED"

	//
	// Flash movies need to have hard coded click URL or clickTAG
	//
	ImageErrorReasonFLASH_WITHOUT_HARD_CODED_CLICK_URL ImageErrorReason = "FLASH_WITHOUT_HARD_CODED_CLICK_URL"

	//
	// Uploaded flash file is corrupted.
	//
	ImageErrorReasonINVALID_FLASH_FILE ImageErrorReason = "INVALID_FLASH_FILE"

	//
	// Uploaded flash file can be parsed, but the click tag can not be fixed properly.
	//
	ImageErrorReasonFAILED_TO_FIX_CLICK_TAG_IN_FLASH ImageErrorReason = "FAILED_TO_FIX_CLICK_TAG_IN_FLASH"

	//
	// Flash movie accesses network resources
	//
	ImageErrorReasonFLASH_ACCESSES_NETWORK_RESOURCES ImageErrorReason = "FLASH_ACCESSES_NETWORK_RESOURCES"

	//
	// Flash movie attempts to call external javascript code
	//
	ImageErrorReasonFLASH_EXTERNAL_JS_CALL ImageErrorReason = "FLASH_EXTERNAL_JS_CALL"

	//
	// Flash movie attempts to call flash system commands
	//
	ImageErrorReasonFLASH_EXTERNAL_FS_CALL ImageErrorReason = "FLASH_EXTERNAL_FS_CALL"

	//
	// Image file is too large.
	//
	ImageErrorReasonFILE_TOO_LARGE ImageErrorReason = "FILE_TOO_LARGE"

	//
	// Image data is too large.
	//
	ImageErrorReasonIMAGE_DATA_TOO_LARGE ImageErrorReason = "IMAGE_DATA_TOO_LARGE"

	//
	// Error while processing the image.
	//
	ImageErrorReasonIMAGE_PROCESSING_ERROR ImageErrorReason = "IMAGE_PROCESSING_ERROR"

	//
	// Image is too small.
	//
	ImageErrorReasonIMAGE_TOO_SMALL ImageErrorReason = "IMAGE_TOO_SMALL"

	//
	// Input was invalid.
	//
	ImageErrorReasonINVALID_INPUT ImageErrorReason = "INVALID_INPUT"

	//
	// There was a problem reading the image file.
	//
	ImageErrorReasonPROBLEM_READING_FILE ImageErrorReason = "PROBLEM_READING_FILE"
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
// The reasons for the target error.
//
type MediaErrorReason string

const (

	//
	// Cannot add a standard icon type
	//
	MediaErrorReasonCANNOT_ADD_STANDARD_ICON MediaErrorReason = "CANNOT_ADD_STANDARD_ICON"

	//
	// May only select Standard Icons alone
	//
	MediaErrorReasonCANNOT_SELECT_STANDARD_ICON_WITH_OTHER_TYPES MediaErrorReason = "CANNOT_SELECT_STANDARD_ICON_WITH_OTHER_TYPES"

	//
	// Image contains both a media ID and media data.
	//
	MediaErrorReasonCANNOT_SPECIFY_MEDIA_ID_AND_DATA MediaErrorReason = "CANNOT_SPECIFY_MEDIA_ID_AND_DATA"

	//
	// A media with given type and reference id already exists
	//
	MediaErrorReasonDUPLICATE_MEDIA MediaErrorReason = "DUPLICATE_MEDIA"

	//
	// A required field was not specified or is an empty string.
	//
	MediaErrorReasonEMPTY_FIELD MediaErrorReason = "EMPTY_FIELD"

	//
	// A media may only be modified once per call
	//
	MediaErrorReasonENTITY_REFERENCED_IN_MULTIPLE_OPS MediaErrorReason = "ENTITY_REFERENCED_IN_MULTIPLE_OPS"

	//
	// Field is not supported for the media sub type.
	//
	MediaErrorReasonFIELD_NOT_SUPPORTED_FOR_MEDIA_SUB_TYPE MediaErrorReason = "FIELD_NOT_SUPPORTED_FOR_MEDIA_SUB_TYPE"

	//
	// The media id is invalid
	//
	MediaErrorReasonINVALID_MEDIA_ID MediaErrorReason = "INVALID_MEDIA_ID"

	//
	// The media subtype is invalid
	//
	MediaErrorReasonINVALID_MEDIA_SUB_TYPE MediaErrorReason = "INVALID_MEDIA_SUB_TYPE"

	//
	// The media type is invalid
	//
	MediaErrorReasonINVALID_MEDIA_TYPE MediaErrorReason = "INVALID_MEDIA_TYPE"

	//
	// The mimetype is invalid
	//
	MediaErrorReasonINVALID_MIME_TYPE MediaErrorReason = "INVALID_MIME_TYPE"

	//
	// The media reference id is invalid
	//
	MediaErrorReasonINVALID_REFERENCE_ID MediaErrorReason = "INVALID_REFERENCE_ID"

	//
	// The YouTube video id is invalid
	//
	MediaErrorReasonINVALID_YOU_TUBE_ID MediaErrorReason = "INVALID_YOU_TUBE_ID"

	//
	// Media has failed transcoding
	//
	MediaErrorReasonMEDIA_FAILED_TRANSCODING MediaErrorReason = "MEDIA_FAILED_TRANSCODING"

	//
	// Media has not been transcoded
	//
	MediaErrorReasonMEDIA_NOT_TRANSCODED MediaErrorReason = "MEDIA_NOT_TRANSCODED"

	//
	// The MediaType does not match the actual media object's type
	//
	MediaErrorReasonMEDIA_TYPE_DOES_NOT_MATCH_OBJECT_TYPE MediaErrorReason = "MEDIA_TYPE_DOES_NOT_MATCH_OBJECT_TYPE"

	//
	// None of the fields have been specified.
	//
	MediaErrorReasonNO_FIELDS_SPECIFIED MediaErrorReason = "NO_FIELDS_SPECIFIED"

	//
	// One of reference Id or media Id must be specified
	//
	MediaErrorReasonNULL_REFERENCE_ID_AND_MEDIA_ID MediaErrorReason = "NULL_REFERENCE_ID_AND_MEDIA_ID"

	//
	// The string has too many characters.
	//
	MediaErrorReasonTOO_LONG MediaErrorReason = "TOO_LONG"

	//
	// The specified operation is not supported.  Only ADD, SET, and REMOVE
	// are supported
	//
	MediaErrorReasonUNSUPPORTED_OPERATION MediaErrorReason = "UNSUPPORTED_OPERATION"

	//
	// The specified type is not supported.
	//
	MediaErrorReasonUNSUPPORTED_TYPE MediaErrorReason = "UNSUPPORTED_TYPE"

	//
	// YouTube is unavailable for requesting video data.
	//
	MediaErrorReasonYOU_TUBE_SERVICE_UNAVAILABLE MediaErrorReason = "YOU_TUBE_SERVICE_UNAVAILABLE"

	//
	// The YouTube video has a non positive duration.
	//
	MediaErrorReasonYOU_TUBE_VIDEO_HAS_NON_POSITIVE_DURATION MediaErrorReason = "YOU_TUBE_VIDEO_HAS_NON_POSITIVE_DURATION"

	//
	// The YouTube video id is syntactically valid but the video was not found.
	//
	MediaErrorReasonYOU_TUBE_VIDEO_NOT_FOUND MediaErrorReason = "YOU_TUBE_VIDEO_NOT_FOUND"
)

//
// Reason for bidding error.
//
type MultiplierErrorReason string

const (

	//
	// Multiplier value is too high
	//
	MultiplierErrorReasonMULTIPLIER_TOO_HIGH MultiplierErrorReason = "MULTIPLIER_TOO_HIGH"

	//
	// Multiplier value is too low
	//
	MultiplierErrorReasonMULTIPLIER_TOO_LOW MultiplierErrorReason = "MULTIPLIER_TOO_LOW"

	//
	// Too many fractional digits
	//
	MultiplierErrorReasonTOO_MANY_FRACTIONAL_DIGITS MultiplierErrorReason = "TOO_MANY_FRACTIONAL_DIGITS"

	//
	// A multiplier cannot be set for this bidding strategy
	//
	MultiplierErrorReasonMULTIPLIER_NOT_ALLOWED_FOR_BIDDING_STRATEGY MultiplierErrorReason = "MULTIPLIER_NOT_ALLOWED_FOR_BIDDING_STRATEGY"

	//
	// A multiplier cannot be set when there is no base bid (e.g., content max cpc)
	//
	MultiplierErrorReasonMULTIPLIER_NOT_ALLOWED_WHEN_BASE_BID_IS_MISSING MultiplierErrorReason = "MULTIPLIER_NOT_ALLOWED_WHEN_BASE_BID_IS_MISSING"

	//
	// A bid multiplier must be specified
	//
	MultiplierErrorReasonNO_MULTIPLIER_SPECIFIED MultiplierErrorReason = "NO_MULTIPLIER_SPECIFIED"

	//
	// Multiplier causes bid to exceed daily budget
	//
	MultiplierErrorReasonMULTIPLIER_CAUSES_BID_TO_EXCEED_DAILY_BUDGET MultiplierErrorReason = "MULTIPLIER_CAUSES_BID_TO_EXCEED_DAILY_BUDGET"

	//
	// Multiplier causes bid to exceed monthly budget
	//
	MultiplierErrorReasonMULTIPLIER_CAUSES_BID_TO_EXCEED_MONTHLY_BUDGET MultiplierErrorReason = "MULTIPLIER_CAUSES_BID_TO_EXCEED_MONTHLY_BUDGET"

	//
	// Multiplier causes bid to exceed custom budget
	//
	MultiplierErrorReasonMULTIPLIER_CAUSES_BID_TO_EXCEED_CUSTOM_BUDGET MultiplierErrorReason = "MULTIPLIER_CAUSES_BID_TO_EXCEED_CUSTOM_BUDGET"

	//
	// Multiplier causes bid to exceed maximum allowed bid
	//
	MultiplierErrorReasonMULTIPLIER_CAUSES_BID_TO_EXCEED_MAX_ALLOWED_BID MultiplierErrorReason = "MULTIPLIER_CAUSES_BID_TO_EXCEED_MAX_ALLOWED_BID"

	//
	// Multiplier causes bid to become less than the minimum bid allowed
	//
	MultiplierErrorReasonBID_LESS_THAN_MIN_ALLOWED_BID_WITH_MULTIPLIER MultiplierErrorReason = "BID_LESS_THAN_MIN_ALLOWED_BID_WITH_MULTIPLIER"

	//
	// Multiplier type (cpc vs. cpm) needs to match campaign's bidding strategy
	//
	MultiplierErrorReasonMULTIPLIER_AND_BIDDING_STRATEGY_TYPE_MISMATCH MultiplierErrorReason = "MULTIPLIER_AND_BIDDING_STRATEGY_TYPE_MISMATCH"

	MultiplierErrorReasonMULTIPLIER_ERROR MultiplierErrorReason = "MULTIPLIER_ERROR"
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
// The reasons for the setting error.
//
type SettingErrorReason string

const (

	//
	// The campaign already has a setting of the type that is being added.
	//
	SettingErrorReasonDUPLICATE_SETTING_TYPE SettingErrorReason = "DUPLICATE_SETTING_TYPE"

	//
	// The campaign setting is not available for this AdWords account.
	//
	SettingErrorReasonSETTING_TYPE_IS_NOT_AVAILABLE SettingErrorReason = "SETTING_TYPE_IS_NOT_AVAILABLE"

	//
	// The setting is not compatible with the campaign.
	//
	SettingErrorReasonSETTING_TYPE_IS_NOT_COMPATIBLE_WITH_CAMPAIGN SettingErrorReason = "SETTING_TYPE_IS_NOT_COMPATIBLE_WITH_CAMPAIGN"

	//
	// The supplied TargetingSetting contains an invalid CriterionTypeGroup. See
	// {@link CriterionTypeGroup} documentation for CriterionTypeGroups allowed in Campaign or
	// AdGroup TargetingSettings.
	//
	SettingErrorReasonTARGETING_SETTING_CONTAINS_INVALID_CRITERION_TYPE_GROUP SettingErrorReason = "TARGETING_SETTING_CONTAINS_INVALID_CRITERION_TYPE_GROUP"

	//
	// Starting with AdWords API v201802, TargetingSetting must not explicitly set any of the
	// Demographic CriterionTypeGroups (AGE_RANGE, GENDER, PARENT, INCOME_RANGE) to false (it's
	// okay to not set them at all, in which case the system will set them to true automatically).
	//
	SettingErrorReasonTARGETING_SETTING_DEMOGRAPHIC_CRITERION_TYPE_GROUPS_MUST_BE_SET_TO_TARGET_ALL SettingErrorReason = "TARGETING_SETTING_DEMOGRAPHIC_CRITERION_TYPE_GROUPS_MUST_BE_SET_TO_TARGET_ALL"

	//
	// Starting with AdWords API v201802,TargetingSetting cannot change any of the Demographic
	// CriterionTypeGroups (AGE_RANGE, GENDER, PARENT, INCOME_RANGE) from true to false.
	//
	SettingErrorReasonTARGETING_SETTING_CANNOT_CHANGE_TARGET_ALL_TO_FALSE_FOR_DEMOGRAPHIC_CRITERION_TYPE_GROUP SettingErrorReason = "TARGETING_SETTING_CANNOT_CHANGE_TARGET_ALL_TO_FALSE_FOR_DEMOGRAPHIC_CRITERION_TYPE_GROUP"

	//
	// At least one feed id should be present.
	//
	SettingErrorReasonDYNAMIC_SEARCH_ADS_SETTING_AT_LEAST_ONE_FEED_ID_MUST_BE_PRESENT SettingErrorReason = "DYNAMIC_SEARCH_ADS_SETTING_AT_LEAST_ONE_FEED_ID_MUST_BE_PRESENT"

	//
	// The supplied DynamicSearchAdsSetting contains an invalid domain name.
	//
	SettingErrorReasonDYNAMIC_SEARCH_ADS_SETTING_CONTAINS_INVALID_DOMAIN_NAME SettingErrorReason = "DYNAMIC_SEARCH_ADS_SETTING_CONTAINS_INVALID_DOMAIN_NAME"

	//
	// The supplied DynamicSearchAdsSetting contains a subdomain name.
	//
	SettingErrorReasonDYNAMIC_SEARCH_ADS_SETTING_CONTAINS_SUBDOMAIN_NAME SettingErrorReason = "DYNAMIC_SEARCH_ADS_SETTING_CONTAINS_SUBDOMAIN_NAME"

	//
	// The supplied DynamicSearchAdsSetting contains an invalid language code.
	//
	SettingErrorReasonDYNAMIC_SEARCH_ADS_SETTING_CONTAINS_INVALID_LANGUAGE_CODE SettingErrorReason = "DYNAMIC_SEARCH_ADS_SETTING_CONTAINS_INVALID_LANGUAGE_CODE"

	//
	// TargetingSettings in search campaigns should not have CriterionTypeGroup.PLACEMENT
	// set to targetAll.
	//
	SettingErrorReasonTARGET_ALL_IS_NOT_ALLOWED_FOR_PLACEMENT_IN_SEARCH_CAMPAIGN SettingErrorReason = "TARGET_ALL_IS_NOT_ALLOWED_FOR_PLACEMENT_IN_SEARCH_CAMPAIGN"

	//
	// Duplicate description in universal app setting description field.
	//
	SettingErrorReasonUNIVERSAL_APP_CAMPAIGN_SETTING_DUPLICATE_DESCRIPTION SettingErrorReason = "UNIVERSAL_APP_CAMPAIGN_SETTING_DUPLICATE_DESCRIPTION"

	//
	// Description line width is too long in universal app setting description field.
	//
	SettingErrorReasonUNIVERSAL_APP_CAMPAIGN_SETTING_DESCRIPTION_LINE_WIDTH_TOO_LONG SettingErrorReason = "UNIVERSAL_APP_CAMPAIGN_SETTING_DESCRIPTION_LINE_WIDTH_TOO_LONG"

	//
	// Universal app setting appId field cannot be modified for COMPLETE campaigns.
	//
	SettingErrorReasonUNIVERSAL_APP_CAMPAIGN_SETTING_APP_ID_CANNOT_BE_MODIFIED SettingErrorReason = "UNIVERSAL_APP_CAMPAIGN_SETTING_APP_ID_CANNOT_BE_MODIFIED"

	//
	// YoutubeVideoMediaIds in universal app setting cannot exceed size limit.
	//
	SettingErrorReasonTOO_MANY_YOUTUBE_MEDIA_IDS_IN_UNIVERSAL_APP_CAMPAIGN SettingErrorReason = "TOO_MANY_YOUTUBE_MEDIA_IDS_IN_UNIVERSAL_APP_CAMPAIGN"

	//
	// ImageMediaIds in universal app setting cannot exceed size limit.
	//
	SettingErrorReasonTOO_MANY_IMAGE_MEDIA_IDS_IN_UNIVERSAL_APP_CAMPAIGN SettingErrorReason = "TOO_MANY_IMAGE_MEDIA_IDS_IN_UNIVERSAL_APP_CAMPAIGN"

	//
	// Media is incompatible for universal app campaign.
	//
	SettingErrorReasonMEDIA_INCOMPATIBLE_FOR_UNIVERSAL_APP_CAMPAIGN SettingErrorReason = "MEDIA_INCOMPATIBLE_FOR_UNIVERSAL_APP_CAMPAIGN"

	//
	// Too many exclamation marks in universal app campaign ad text ideas.
	//
	SettingErrorReasonTOO_MANY_EXCLAMATION_MARKS SettingErrorReason = "TOO_MANY_EXCLAMATION_MARKS"

	//
	// Unspecified campaign setting error.
	// <span class="constraint Rejected">Used for return value only. An enumeration could not be processed, typically due to incompatibility with your WSDL version.</span>
	//
	SettingErrorReasonUNKNOWN SettingErrorReason = "UNKNOWN"
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
// The reasons for the url error.
//
type UrlErrorReason string

const (

	//
	// The tracking url template is invalid.
	//
	UrlErrorReasonINVALID_TRACKING_URL_TEMPLATE UrlErrorReason = "INVALID_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template contains invalid tag.
	//
	UrlErrorReasonINVALID_TAG_IN_TRACKING_URL_TEMPLATE UrlErrorReason = "INVALID_TAG_IN_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template must contain at least one tag (e.g. {lpurl}),
	// This applies only to tracking url template associated with website ads or product ads.
	//
	UrlErrorReasonMISSING_TRACKING_URL_TEMPLATE_TAG UrlErrorReason = "MISSING_TRACKING_URL_TEMPLATE_TAG"

	//
	// The tracking url template must start with a valid protocol (or lpurl tag).
	//
	UrlErrorReasonMISSING_PROTOCOL_IN_TRACKING_URL_TEMPLATE UrlErrorReason = "MISSING_PROTOCOL_IN_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template starts with an invalid protocol.
	//
	UrlErrorReasonINVALID_PROTOCOL_IN_TRACKING_URL_TEMPLATE UrlErrorReason = "INVALID_PROTOCOL_IN_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template  contains illegal characters.
	//
	UrlErrorReasonMALFORMED_TRACKING_URL_TEMPLATE UrlErrorReason = "MALFORMED_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template must contain a host name (or lpurl tag).
	//
	UrlErrorReasonMISSING_HOST_IN_TRACKING_URL_TEMPLATE UrlErrorReason = "MISSING_HOST_IN_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template has an invalid or missing top level domain extension.
	//
	UrlErrorReasonINVALID_TLD_IN_TRACKING_URL_TEMPLATE UrlErrorReason = "INVALID_TLD_IN_TRACKING_URL_TEMPLATE"

	//
	// The tracking url template contains nested occurrences of the same conditional tag
	// (i.e. {ifmobile:{ifmobile:x}}).
	//
	UrlErrorReasonREDUNDANT_NESTED_TRACKING_URL_TEMPLATE_TAG UrlErrorReason = "REDUNDANT_NESTED_TRACKING_URL_TEMPLATE_TAG"

	//
	// The final url is invalid.
	//
	UrlErrorReasonINVALID_FINAL_URL UrlErrorReason = "INVALID_FINAL_URL"

	//
	// The final url contains invalid tag.
	//
	UrlErrorReasonINVALID_TAG_IN_FINAL_URL UrlErrorReason = "INVALID_TAG_IN_FINAL_URL"

	//
	// The final url contains nested occurrences of the same conditional tag
	// (i.e. {ifmobile:{ifmobile:x}}).
	//
	UrlErrorReasonREDUNDANT_NESTED_FINAL_URL_TAG UrlErrorReason = "REDUNDANT_NESTED_FINAL_URL_TAG"

	//
	// The final url must start with a valid protocol.
	//
	UrlErrorReasonMISSING_PROTOCOL_IN_FINAL_URL UrlErrorReason = "MISSING_PROTOCOL_IN_FINAL_URL"

	//
	// The final url starts with an invalid protocol.
	//
	UrlErrorReasonINVALID_PROTOCOL_IN_FINAL_URL UrlErrorReason = "INVALID_PROTOCOL_IN_FINAL_URL"

	//
	// The final url  contains illegal characters.
	//
	UrlErrorReasonMALFORMED_FINAL_URL UrlErrorReason = "MALFORMED_FINAL_URL"

	//
	// The final url must contain a host name.
	//
	UrlErrorReasonMISSING_HOST_IN_FINAL_URL UrlErrorReason = "MISSING_HOST_IN_FINAL_URL"

	//
	// The tracking url template has an invalid or missing top level domain extension.
	//
	UrlErrorReasonINVALID_TLD_IN_FINAL_URL UrlErrorReason = "INVALID_TLD_IN_FINAL_URL"

	//
	// The final mobile url is invalid.
	//
	UrlErrorReasonINVALID_FINAL_MOBILE_URL UrlErrorReason = "INVALID_FINAL_MOBILE_URL"

	//
	// The final mobile url contains invalid tag.
	//
	UrlErrorReasonINVALID_TAG_IN_FINAL_MOBILE_URL UrlErrorReason = "INVALID_TAG_IN_FINAL_MOBILE_URL"

	//
	// The final mobile url contains nested occurrences of the same conditional tag
	// (i.e. {ifmobile:{ifmobile:x}}).
	//
	UrlErrorReasonREDUNDANT_NESTED_FINAL_MOBILE_URL_TAG UrlErrorReason = "REDUNDANT_NESTED_FINAL_MOBILE_URL_TAG"

	//
	// The final mobile url must start with a valid protocol.
	//
	UrlErrorReasonMISSING_PROTOCOL_IN_FINAL_MOBILE_URL UrlErrorReason = "MISSING_PROTOCOL_IN_FINAL_MOBILE_URL"

	//
	// The final mobile url starts with an invalid protocol.
	//
	UrlErrorReasonINVALID_PROTOCOL_IN_FINAL_MOBILE_URL UrlErrorReason = "INVALID_PROTOCOL_IN_FINAL_MOBILE_URL"

	//
	// The final mobile url  contains illegal characters.
	//
	UrlErrorReasonMALFORMED_FINAL_MOBILE_URL UrlErrorReason = "MALFORMED_FINAL_MOBILE_URL"

	//
	// The final mobile url must contain a host name.
	//
	UrlErrorReasonMISSING_HOST_IN_FINAL_MOBILE_URL UrlErrorReason = "MISSING_HOST_IN_FINAL_MOBILE_URL"

	//
	// The tracking url template has an invalid or missing top level domain extension.
	//
	UrlErrorReasonINVALID_TLD_IN_FINAL_MOBILE_URL UrlErrorReason = "INVALID_TLD_IN_FINAL_MOBILE_URL"

	//
	// The final app url is invalid.
	//
	UrlErrorReasonINVALID_FINAL_APP_URL UrlErrorReason = "INVALID_FINAL_APP_URL"

	//
	// The final app url contains invalid tag.
	//
	UrlErrorReasonINVALID_TAG_IN_FINAL_APP_URL UrlErrorReason = "INVALID_TAG_IN_FINAL_APP_URL"

	//
	// The final app url contains nested occurrences of the same conditional tag
	// (i.e. {ifmobile:{ifmobile:x}}).
	//
	UrlErrorReasonREDUNDANT_NESTED_FINAL_APP_URL_TAG UrlErrorReason = "REDUNDANT_NESTED_FINAL_APP_URL_TAG"

	//
	// More than one app url found for the same OS type.
	//
	UrlErrorReasonMULTIPLE_APP_URLS_FOR_OSTYPE UrlErrorReason = "MULTIPLE_APP_URLS_FOR_OSTYPE"

	//
	// The OS type given for an app url is not valid.
	//
	UrlErrorReasonINVALID_OSTYPE UrlErrorReason = "INVALID_OSTYPE"

	//
	// The protocol given for an app url is not valid. (E.g. "android-app://")
	//
	UrlErrorReasonINVALID_PROTOCOL_FOR_APP_URL UrlErrorReason = "INVALID_PROTOCOL_FOR_APP_URL"

	//
	// The package id (app id) given for an app url is not valid.
	//
	UrlErrorReasonINVALID_PACKAGE_ID_FOR_APP_URL UrlErrorReason = "INVALID_PACKAGE_ID_FOR_APP_URL"

	//
	// The number of url custom parameters for an entity exceeds the maximum limit allowed.
	//
	UrlErrorReasonURL_CUSTOM_PARAMETERS_COUNT_EXCEEDS_LIMIT UrlErrorReason = "URL_CUSTOM_PARAMETERS_COUNT_EXCEEDS_LIMIT"

	//
	// The parameter has isRemove set to true but a value that is non-null.
	//
	UrlErrorReasonURL_CUSTOM_PARAMETER_REMOVAL_WITH_NON_NULL_VALUE UrlErrorReason = "URL_CUSTOM_PARAMETER_REMOVAL_WITH_NON_NULL_VALUE"

	//
	// For add operations, there will not be any existing parameters to delete.
	//
	UrlErrorReasonCANNOT_REMOVE_URL_CUSTOM_PARAMETER_IN_ADD_OPERATION UrlErrorReason = "CANNOT_REMOVE_URL_CUSTOM_PARAMETER_IN_ADD_OPERATION"

	//
	// When the doReplace flag is set to true, individual parameters cannot be deleted.
	//
	UrlErrorReasonCANNOT_REMOVE_URL_CUSTOM_PARAMETER_DURING_FULL_REPLACEMENT UrlErrorReason = "CANNOT_REMOVE_URL_CUSTOM_PARAMETER_DURING_FULL_REPLACEMENT"

	//
	// For ADD operations and when the doReplace flag is set to true, custom parameter values
	// cannot be null.
	//
	UrlErrorReasonNULL_CUSTOM_PARAMETER_VALUE_DURING_ADD_OR_FULL_REPLACEMENT UrlErrorReason = "NULL_CUSTOM_PARAMETER_VALUE_DURING_ADD_OR_FULL_REPLACEMENT"

	//
	// An invalid character appears in the parameter key.
	//
	UrlErrorReasonINVALID_CHARACTERS_IN_URL_CUSTOM_PARAMETER_KEY UrlErrorReason = "INVALID_CHARACTERS_IN_URL_CUSTOM_PARAMETER_KEY"

	//
	// An invalid character appears in the parameter value.
	//
	UrlErrorReasonINVALID_CHARACTERS_IN_URL_CUSTOM_PARAMETER_VALUE UrlErrorReason = "INVALID_CHARACTERS_IN_URL_CUSTOM_PARAMETER_VALUE"

	//
	// The url custom parameter value fails url tag validation.
	//
	UrlErrorReasonINVALID_TAG_IN_URL_CUSTOM_PARAMETER_VALUE UrlErrorReason = "INVALID_TAG_IN_URL_CUSTOM_PARAMETER_VALUE"

	//
	// The custom parameter contains nested occurrences of the same conditional tag
	// (i.e. {ifmobile:{ifmobile:x}}).
	//
	UrlErrorReasonREDUNDANT_NESTED_URL_CUSTOM_PARAMETER_TAG UrlErrorReason = "REDUNDANT_NESTED_URL_CUSTOM_PARAMETER_TAG"

	//
	// The protocol (http:// or https://) is missing.
	//
	UrlErrorReasonMISSING_PROTOCOL UrlErrorReason = "MISSING_PROTOCOL"

	//
	// The url is invalid.
	//
	UrlErrorReasonINVALID_URL UrlErrorReason = "INVALID_URL"

	//
	// Destination Url is deprecated.
	//
	UrlErrorReasonDESTINATION_URL_DEPRECATED UrlErrorReason = "DESTINATION_URL_DEPRECATED"

	//
	// The url contains invalid tag.
	//
	UrlErrorReasonINVALID_TAG_IN_URL UrlErrorReason = "INVALID_TAG_IN_URL"

	//
	// The url must contain at least one tag (e.g. {lpurl}),
	// This applies only to urls associated with website ads or product ads.
	//
	UrlErrorReasonMISSING_URL_TAG UrlErrorReason = "MISSING_URL_TAG"

	UrlErrorReasonDUPLICATE_URL_ID UrlErrorReason = "DUPLICATE_URL_ID"

	UrlErrorReasonINVALID_URL_ID UrlErrorReason = "INVALID_URL_ID"

	UrlErrorReasonURL_ERROR UrlErrorReason = "URL_ERROR"
)

type VideoErrorReason string

const (

	//
	// Invalid video.
	//
	VideoErrorReasonINVALID_VIDEO VideoErrorReason = "INVALID_VIDEO"

	//
	// Storage error.
	//
	VideoErrorReasonSTORAGE_ERROR VideoErrorReason = "STORAGE_ERROR"

	//
	// Bad request.
	//
	VideoErrorReasonBAD_REQUEST VideoErrorReason = "BAD_REQUEST"

	//
	// Server error.
	//
	VideoErrorReasonERROR_GENERATING_STREAMING_URL VideoErrorReason = "ERROR_GENERATING_STREAMING_URL"

	//
	// Unexpected size.
	//
	VideoErrorReasonUNEXPECTED_SIZE VideoErrorReason = "UNEXPECTED_SIZE"

	//
	// Server error.
	//
	VideoErrorReasonSERVER_ERROR VideoErrorReason = "SERVER_ERROR"

	//
	// File too large.
	//
	VideoErrorReasonFILE_TOO_LARGE VideoErrorReason = "FILE_TOO_LARGE"

	//
	// Video processing error.
	//
	VideoErrorReasonVIDEO_PROCESSING_ERROR VideoErrorReason = "VIDEO_PROCESSING_ERROR"

	//
	// Invalid input.
	//
	VideoErrorReasonINVALID_INPUT VideoErrorReason = "INVALID_INPUT"

	//
	// Problem reading file.
	//
	VideoErrorReasonPROBLEM_READING_FILE VideoErrorReason = "PROBLEM_READING_FILE"

	//
	// Invalid ISCI.
	//
	VideoErrorReasonINVALID_ISCI VideoErrorReason = "INVALID_ISCI"

	//
	// Invalid AD-ID.
	//
	VideoErrorReasonINVALID_AD_ID VideoErrorReason = "INVALID_AD_ID"
)

type Get struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 get"`

	//
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Selector *Selector `xml:"selector,omitempty"`
}

type GetResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getResponse"`

	Rval *DraftAsyncErrorPage `xml:"rval,omitempty"`
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

	Rval *DraftAsyncErrorPage `xml:"rval,omitempty"`
}

type AdError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AdErrorReason `xml:"reason,omitempty"`
}

type AdGroupAdError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupAdError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AdGroupAdErrorReason `xml:"reason,omitempty"`
}

type AdGroupCriterionError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupCriterionError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AdGroupCriterionErrorReason `xml:"reason,omitempty"`
}

type AdGroupFeedError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupFeedError"`

	*ApiError

	//
	// Error reason.
	//
	Reason *AdGroupFeedErrorReason `xml:"reason,omitempty"`
}

type AdGroupServiceError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdGroupServiceError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AdGroupServiceErrorReason `xml:"reason,omitempty"`
}

type AdSharingError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AdSharingError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AdSharingErrorReason `xml:"reason,omitempty"`

	SharedAdError *ApiError `xml:"sharedAdError,omitempty"`
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

type BiddingErrors struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 BiddingErrors"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *BiddingErrorsReason `xml:"reason,omitempty"`
}

type CampaignBidModifierError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignBidModifierError"`

	*ApiError

	Reason *CampaignBidModifierErrorReason `xml:"reason,omitempty"`
}

type CampaignCriterionError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignCriterionError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *CampaignCriterionErrorReason `xml:"reason,omitempty"`
}

type CampaignError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *CampaignErrorReason `xml:"reason,omitempty"`
}

type CampaignFeedError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignFeedError"`

	*ApiError

	//
	// Error reason.
	//
	Reason *CampaignFeedErrorReason `xml:"reason,omitempty"`
}

type CampaignPreferenceError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignPreferenceError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *CampaignPreferenceErrorReason `xml:"reason,omitempty"`
}

type CampaignSharedSetError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CampaignSharedSetError"`

	*ApiError

	Reason *CampaignSharedSetErrorReason `xml:"reason,omitempty"`
}

type ClientTermsError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ClientTermsError"`

	*ApiError

	Reason *ClientTermsErrorReason `xml:"reason,omitempty"`
}

type CriterionError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 CriterionError"`

	*ApiError

	Reason *CriterionErrorReason `xml:"reason,omitempty"`
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

type DateRangeError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DateRangeError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DateRangeErrorReason `xml:"reason,omitempty"`
}

type DistinctError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DistinctError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DistinctErrorReason `xml:"reason,omitempty"`
}

type DraftAsyncError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DraftAsyncError"`

	//
	// The error occurred during promotion while updating this Campaign or an entity in this Campaign.
	// This field can only be used with Predicate Operators EQUALS and IN. When using a Predicate
	// with this field, also include a Predicate for the field DraftId.
	// <span class="constraint Selectable">This field can be selected using the value "BaseCampaignId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	BaseCampaignId int64 `xml:"baseCampaignId,omitempty"`

	//
	// The draft that was attempted to be promoted.  This field can only be used with Predicate
	// Operators EQUALS and IN. When using a Predicate with this field, also include a Predicate for
	// the field BaseCampaignId.
	// <span class="constraint Selectable">This field can be selected using the value "DraftId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DraftId int64 `xml:"draftId,omitempty"`

	//
	// The draft Campaign that was attempted to be promoted.
	// <span class="constraint Selectable">This field can be selected using the value "DraftCampaignId".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DraftCampaignId int64 `xml:"draftCampaignId,omitempty"`

	//
	// The error that occurred while promoting the draft.
	// <span class="constraint Selectable">This field can be selected using the value "AsyncError".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	AsyncError *ApiError `xml:"asyncError,omitempty"`

	//
	// The error occurred during promotion while updating this AdGroup or an entity in this AdGroup.
	// <span class="constraint Selectable">This field can be selected using the value "BaseAdGroupId".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	BaseAdGroupId int64 `xml:"baseAdGroupId,omitempty"`

	//
	// The draft AdGroup that was attempted to be promoted.
	// <span class="constraint Selectable">This field can be selected using the value "DraftAdGroupId".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	DraftAdGroupId int64 `xml:"draftAdGroupId,omitempty"`
}

type DraftAsyncErrorPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DraftAsyncErrorPage"`

	*Page

	Entries []*DraftAsyncError `xml:"entries,omitempty"`
}

type DraftError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DraftError"`

	*ApiError

	Reason *DraftErrorReason `xml:"reason,omitempty"`
}

type EntityAccessDenied struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 EntityAccessDenied"`

	*ApiError

	//
	// Reason for this error.
	//
	Reason *EntityAccessDeniedReason `xml:"reason,omitempty"`
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

type FeedError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 FeedError"`

	*ApiError

	//
	// The cause of the error.
	//
	Reason *FeedErrorReason `xml:"reason,omitempty"`
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

type FunctionError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 FunctionError"`

	*ApiError

	//
	// The error reason represented by an enum
	//
	Reason *FunctionErrorReason `xml:"reason,omitempty"`
}

type IdError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 IdError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *IdErrorReason `xml:"reason,omitempty"`
}

type ImageError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 ImageError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *ImageErrorReason `xml:"reason,omitempty"`
}

type InternalApiError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 InternalApiError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *InternalApiErrorReason `xml:"reason,omitempty"`
}

type MediaError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MediaError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *MediaErrorReason `xml:"reason,omitempty"`
}

type MultiplierError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MultiplierError"`

	*ApiError

	Reason *MultiplierErrorReason `xml:"reason,omitempty"`
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

type PolicyViolationError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 PolicyViolationError"`

	*ApiError

	//
	// Unique identifier for the violation.
	//
	Key *PolicyViolationKey `xml:"key,omitempty"`

	//
	// Name of policy suitable for display to users. In the user's preferred
	// language.
	//
	ExternalPolicyName string `xml:"externalPolicyName,omitempty"`

	//
	// Url with writeup about the policy.
	//
	ExternalPolicyUrl string `xml:"externalPolicyUrl,omitempty"`

	//
	// Localized description of the violation.
	//
	ExternalPolicyDescription string `xml:"externalPolicyDescription,omitempty"`

	//
	// Whether user can file an exemption request for this violation.
	//
	IsExemptable bool `xml:"isExemptable,omitempty"`

	//
	// Lists the parts that violate the policy.
	//
	ViolatingParts []*PolicyViolationErrorPart `xml:"violatingParts,omitempty"`
}

type PolicyViolationErrorPart struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 PolicyViolationError.Part"`

	//
	// Index of the starting position of the violating text within the line.
	//
	Index int32 `xml:"index,omitempty"`

	//
	// The length of the violating text.
	//
	Length int32 `xml:"length,omitempty"`
}

type PolicyViolationKey struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 PolicyViolationKey"`

	//
	// Unique id of the violated policy.
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	PolicyName string `xml:"policyName,omitempty"`

	//
	// The text that violates the policy if specified. Otherwise, refers to the
	// policy in general (e.g. when requesting to be exempt from the whole
	// policy).
	//
	// May be null for criterion exemptions, in which case this refers to the
	// whole policy. Must be specified for ad exemptions.
	//
	ViolatingText string `xml:"violatingText,omitempty"`
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

type SettingError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 SettingError"`

	*ApiError

	//
	// The setting error reason.
	//
	Reason *SettingErrorReason `xml:"reason,omitempty"`
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

type UrlError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 UrlError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *UrlErrorReason `xml:"reason,omitempty"`
}

type VideoError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 VideoError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *VideoErrorReason `xml:"reason,omitempty"`
}

type DraftAsyncErrorServiceInterface struct {
	client *SOAPClient
}

func NewDraftAsyncErrorServiceInterface(url string, tls bool, auth *BasicAuth) *DraftAsyncErrorServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &DraftAsyncErrorServiceInterface{
		client: client,
	}
}

func NewDraftAsyncErrorServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *DraftAsyncErrorServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &DraftAsyncErrorServiceInterface{
		client: client,
	}
}

func (service *DraftAsyncErrorServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *DraftAsyncErrorServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Returns a DraftAsyncErrorPage that contains a list of DraftAsyncErrors matching the selector.

   @throws {#link com.google.ads.api.services.common.error.ApiException} if problems occurred
   while retrieving the results.
*/
func (service *DraftAsyncErrorServiceInterface) Get(request *Get) (*GetResponse, error) {
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
   Returns a DraftAsyncErrorPage that contains a list of DraftAsyncErrors matching the query.

   @throws {#link com.google.ads.api.services.common.error.ApiException} if problems occurred
   while retrieving the results.
*/
func (service *DraftAsyncErrorServiceInterface) Query(request *Query) (*QueryResponse, error) {
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
