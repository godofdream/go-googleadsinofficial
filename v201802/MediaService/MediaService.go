package MediaService

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

type AudioErrorReason string

const (
	AudioErrorReasonINVALID_AUDIO AudioErrorReason = "INVALID_AUDIO"

	AudioErrorReasonPROBLEM_READING_AUDIO_FILE AudioErrorReason = "PROBLEM_READING_AUDIO_FILE"

	AudioErrorReasonERROR_STORING_AUDIO AudioErrorReason = "ERROR_STORING_AUDIO"

	AudioErrorReasonFILE_TOO_LARGE AudioErrorReason = "FILE_TOO_LARGE"

	AudioErrorReasonUNSUPPORTED_AUDIO AudioErrorReason = "UNSUPPORTED_AUDIO"

	AudioErrorReasonERROR_GENERATING_STREAMING_URL AudioErrorReason = "ERROR_GENERATING_STREAMING_URL"
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
// Media types
//
type MediaMediaType string

const (

	//
	// Audio file.
	//
	MediaMediaTypeAUDIO MediaMediaType = "AUDIO"

	//
	// Animated image, such as animated GIF.
	//
	MediaMediaTypeDYNAMIC_IMAGE MediaMediaType = "DYNAMIC_IMAGE"

	//
	// Small image; used for map ad.
	//
	MediaMediaTypeICON MediaMediaType = "ICON"

	//
	// Static image; for image ad.
	//
	MediaMediaTypeIMAGE MediaMediaType = "IMAGE"

	//
	// Predefined standard icon; used for map ads.
	//
	MediaMediaTypeSTANDARD_ICON MediaMediaType = "STANDARD_ICON"

	//
	// Video file.
	//
	MediaMediaTypeVIDEO MediaMediaType = "VIDEO"

	//
	// ZIP file; used in fields of template ads.
	//
	MediaMediaTypeMEDIA_BUNDLE MediaMediaType = "MEDIA_BUNDLE"
)

//
// Mime types
//
type MediaMimeType string

const (

	//
	// MIME type of image/jpeg
	//
	MediaMimeTypeIMAGE_JPEG MediaMimeType = "IMAGE_JPEG"

	//
	// MIME type of image/gif
	//
	MediaMimeTypeIMAGE_GIF MediaMimeType = "IMAGE_GIF"

	//
	// MIME type of image/png
	//
	MediaMimeTypeIMAGE_PNG MediaMimeType = "IMAGE_PNG"

	//
	// MIME type of application/x-shockwave-flash
	//
	MediaMimeTypeFLASH MediaMimeType = "FLASH"

	//
	// MIME type of text/html
	//
	MediaMimeTypeTEXT_HTML MediaMimeType = "TEXT_HTML"

	//
	// MIME type of application/pdf
	//
	MediaMimeTypePDF MediaMimeType = "PDF"

	//
	// MIME type of application/msword
	//
	MediaMimeTypeMSWORD MediaMimeType = "MSWORD"

	//
	// MIME type of application/vnd.ms-excel
	//
	MediaMimeTypeMSEXCEL MediaMimeType = "MSEXCEL"

	//
	// MIME type of application/rtf
	//
	MediaMimeTypeRTF MediaMimeType = "RTF"

	//
	// MIME type of audio/wav
	//
	MediaMimeTypeAUDIO_WAV MediaMimeType = "AUDIO_WAV"

	//
	// MIME type of audio/mp3
	//
	MediaMimeTypeAUDIO_MP3 MediaMimeType = "AUDIO_MP3"

	//
	// MIME type of application/x-html5-ad-zip
	//
	MediaMimeTypeHTML5_AD_ZIP MediaMimeType = "HTML5_AD_ZIP"
)

//
// Sizes for retrieving the original media
//
type MediaSize string

const (

	//
	// Full size of Media.
	//
	MediaSizeFULL MediaSize = "FULL"

	//
	// Shunken size of media.
	//
	MediaSizeSHRUNKEN MediaSize = "SHRUNKEN"

	//
	// Preview size of media.
	//
	MediaSizePREVIEW MediaSize = "PREVIEW"

	//
	// Video thumbnail size of Media.
	//
	MediaSizeVIDEO_THUMBNAIL MediaSize = "VIDEO_THUMBNAIL"
)

//
// Enumeration of the reasons for the {@link MediaBundleError}
//
type MediaBundleErrorReason string

const (

	//
	// The entryPoint field cannot be set using the <code>MediaService</code>.
	//
	MediaBundleErrorReasonENTRY_POINT_CANNOT_BE_SET_USING_MEDIA_SERVICE MediaBundleErrorReason = "ENTRY_POINT_CANNOT_BE_SET_USING_MEDIA_SERVICE"

	//
	// There was a problem with the request.
	//
	MediaBundleErrorReasonBAD_REQUEST MediaBundleErrorReason = "BAD_REQUEST"

	//
	// HTML5 ads using DoubleClick Studio created ZIP files are not supported.
	//
	MediaBundleErrorReasonDOUBLECLICK_BUNDLE_NOT_ALLOWED MediaBundleErrorReason = "DOUBLECLICK_BUNDLE_NOT_ALLOWED"

	//
	// Cannot reference URL external to the media bundle.
	//
	MediaBundleErrorReasonEXTERNAL_URL_NOT_ALLOWED MediaBundleErrorReason = "EXTERNAL_URL_NOT_ALLOWED"

	//
	// Media bundle file is too large.
	//
	MediaBundleErrorReasonFILE_TOO_LARGE MediaBundleErrorReason = "FILE_TOO_LARGE"

	//
	// ZIP file from Google Web Designer is not published.
	//
	MediaBundleErrorReasonGOOGLE_WEB_DESIGNER_ZIP_FILE_NOT_PUBLISHED MediaBundleErrorReason = "GOOGLE_WEB_DESIGNER_ZIP_FILE_NOT_PUBLISHED"

	//
	// Input was invalid.
	//
	MediaBundleErrorReasonINVALID_INPUT MediaBundleErrorReason = "INVALID_INPUT"

	//
	// There was a problem with the media bundle.
	//
	MediaBundleErrorReasonINVALID_MEDIA_BUNDLE MediaBundleErrorReason = "INVALID_MEDIA_BUNDLE"

	//
	// There was a problem with one or more of the media bundle entries.
	//
	MediaBundleErrorReasonINVALID_MEDIA_BUNDLE_ENTRY MediaBundleErrorReason = "INVALID_MEDIA_BUNDLE_ENTRY"

	//
	// The media bundle contains a file with an unknown mime type
	//
	MediaBundleErrorReasonINVALID_MIME_TYPE MediaBundleErrorReason = "INVALID_MIME_TYPE"

	//
	// The media bundle contain an invalid asset path.
	//
	MediaBundleErrorReasonINVALID_PATH MediaBundleErrorReason = "INVALID_PATH"

	//
	// HTML5 ad is trying to reference an asset not in .ZIP file
	//
	MediaBundleErrorReasonINVALID_URL_REFERENCE MediaBundleErrorReason = "INVALID_URL_REFERENCE"

	//
	// Media data is too large.
	//
	MediaBundleErrorReasonMEDIA_DATA_TOO_LARGE MediaBundleErrorReason = "MEDIA_DATA_TOO_LARGE"

	//
	// The media bundle contains no primary entry.
	//
	MediaBundleErrorReasonMISSING_PRIMARY_MEDIA_BUNDLE_ENTRY MediaBundleErrorReason = "MISSING_PRIMARY_MEDIA_BUNDLE_ENTRY"

	//
	// There was an error on the server.
	//
	MediaBundleErrorReasonSERVER_ERROR MediaBundleErrorReason = "SERVER_ERROR"

	//
	// The image could not be stored.
	//
	MediaBundleErrorReasonSTORAGE_ERROR MediaBundleErrorReason = "STORAGE_ERROR"

	//
	// Media bundle created with the Swiffy tool is not allowed.
	//
	MediaBundleErrorReasonSWIFFY_BUNDLE_NOT_ALLOWED MediaBundleErrorReason = "SWIFFY_BUNDLE_NOT_ALLOWED"

	//
	// The media bundle contains too many files.
	//
	MediaBundleErrorReasonTOO_MANY_FILES MediaBundleErrorReason = "TOO_MANY_FILES"

	//
	// The media bundle is not of legal dimensions.
	//
	MediaBundleErrorReasonUNEXPECTED_SIZE MediaBundleErrorReason = "UNEXPECTED_SIZE"

	//
	// Google Web Designer not created for "AdWords" environment.
	//
	MediaBundleErrorReasonUNSUPPORTED_GOOGLE_WEB_DESIGNER_ENVIRONMENT MediaBundleErrorReason = "UNSUPPORTED_GOOGLE_WEB_DESIGNER_ENVIRONMENT"

	//
	// Unsupported HTML5 feature in HTML5 asset.
	//
	MediaBundleErrorReasonUNSUPPORTED_HTML5_FEATURE MediaBundleErrorReason = "UNSUPPORTED_HTML5_FEATURE"

	//
	// URL in HTML5 entry is not ssl compliant.
	//
	MediaBundleErrorReasonURL_IN_MEDIA_BUNDLE_NOT_SSL_COMPLIANT MediaBundleErrorReason = "URL_IN_MEDIA_BUNDLE_NOT_SSL_COMPLIANT"
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
	ServiceSelector *Selector `xml:"serviceSelector,omitempty"`
}

type GetResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 getResponse"`

	Rval *MediaPage `xml:"rval,omitempty"`
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

	Rval *MediaPage `xml:"rval,omitempty"`
}

type Upload struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 upload"`

	Media []*Media `xml:"media,omitempty"`
}

type UploadResponse struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 uploadResponse"`

	Rval []*Media `xml:"rval,omitempty"`
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

type Audio struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Audio"`

	*Media

	//
	// The duration of the associated audio, in milliseconds.
	// <span class="constraint Selectable">This field can be selected using the value "DurationMillis".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	DurationMillis int64 `xml:"durationMillis,omitempty"`

	//
	// The streaming URL of the audio.
	// <span class="constraint Selectable">This field can be selected using the value "StreamingUrl".</span>
	//
	StreamingUrl string `xml:"streamingUrl,omitempty"`

	//
	// Indicates whether the audio is ready to play on the web.
	// <span class="constraint Selectable">This field can be selected using the value "ReadyToPlayOnTheWeb".</span>
	//
	ReadyToPlayOnTheWeb bool `xml:"readyToPlayOnTheWeb,omitempty"`
}

type AudioError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 AudioError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *AudioErrorReason `xml:"reason,omitempty"`
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

type Dimensions struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Dimensions"`

	//
	// Width of the dimension
	// <span class="constraint Selectable">This field can be selected using the value "Width".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Width int32 `xml:"width,omitempty"`

	//
	// Height of the dimension
	// <span class="constraint Selectable">This field can be selected using the value "Height".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null}.</span>
	//
	Height int32 `xml:"height,omitempty"`
}

type DistinctError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 DistinctError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *DistinctErrorReason `xml:"reason,omitempty"`
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

type Image struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Image"`

	*Media

	//
	// Raw image data.
	//
	Data []byte `xml:"data,omitempty"`
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

type Media struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Media"`

	//
	// ID of this media object.
	// <span class="constraint Selectable">This field can be selected using the value "MediaId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint Required">This field is required and should not be {@code null} when it is contained within {@link Operator}s : SET, REMOVE.</span>
	//
	MediaId int64 `xml:"mediaId,omitempty"`

	//
	// Type of media object. Required when using {@link MediaService#upload} to upload a new media
	// file. MEDIA_BUNDLE, ICON, IMAGE, and DYNAMIC_IMAGE are the supported MediaTypes to upload.
	// <span class="constraint Selectable">This field can be selected using the value "Type".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE and SET.</span>
	//
	Type_ *MediaMediaType `xml:"type,omitempty"`

	//
	// Media reference ID key.
	// <span class="constraint Selectable">This field can be selected using the value "ReferenceId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE and SET.</span>
	//
	ReferenceId int64 `xml:"referenceId,omitempty"`

	//
	// Various dimension sizes for the media. Only applies to image media (and video media for
	// video thumbnails).
	// <span class="constraint Selectable">This field can be selected using the value "Dimensions".</span>
	//
	Dimensions []*Media_Size_DimensionsMapEntry `xml:"dimensions,omitempty"`

	//
	// URLs pointing to the resized media for the given sizes. Only applies to image media.
	// <span class="constraint Selectable">This field can be selected using the value "Urls".</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	Urls []*Media_Size_StringMapEntry `xml:"urls,omitempty"`

	//
	// The mime type of the media.
	// <span class="constraint Selectable">This field can be selected using the value "MimeType".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE and SET.</span>
	//
	MimeType *MediaMimeType `xml:"mimeType,omitempty"`

	//
	// The URL of where the original media was downloaded from (or a file name).
	// <span class="constraint Selectable">This field can be selected using the value "SourceUrl".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE and SET.</span>
	//
	SourceUrl string `xml:"sourceUrl,omitempty"`

	//
	// The name of the media. The name can be used by clients to
	// help identify previously uploaded media.
	// <span class="constraint Selectable">This field can be selected using the value "Name".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	Name string `xml:"name,omitempty"`

	//
	// The size of the media file in bytes.
	// <span class="constraint Selectable">This field can be selected using the value "FileSize".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE and SET.</span>
	//
	FileSize int64 `xml:"fileSize,omitempty"`

	//
	// Media creation date in the format YYYY-MM-DD HH:MM:SS+TZ.
	// This is not updatable and not specifiable.
	// <span class="constraint Selectable">This field can be selected using the value "CreationTime".</span><span class="constraint Filterable">This field can be filtered on.</span>
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API for the following {@link Operator}s: REMOVE and SET.</span>
	//
	CreationTime string `xml:"creationTime,omitempty"`

	//
	// Indicates that this instance is a subtype of Media.
	// Although this field is returned in the response, it is ignored on input
	// and cannot be selected. Specify xsi:type instead.
	//
	MediaType string `xml:"Media.Type,omitempty"`
}

type MediaBundle struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MediaBundle"`

	*Media

	//
	// Raw zipped data.
	//
	Data []byte `xml:"data,omitempty"`

	//
	// URL pointing to the data for the MediaBundle data.
	// <span class="constraint ReadOnly">This field is read only and will be ignored when sent to the API.</span>
	//
	MediaBundleUrl string `xml:"mediaBundleUrl,omitempty"`

	//
	// Entry in the ZIP archive used to display the <code>MediaBundle</code> in an
	// <code>Ad</code>. This field can only be set and returned when the <code>MediaBundle</code> is
	// used with the <code>AdGroupAdService</code>. If this field is set when calling
	// <code>MediaService</code>, an error will be returned.
	//
	// <p>To use a <code>MediaBundle</code> that was created with the <code>MediaService</code> in
	// an <code>Ad</code>, create a bundle and set the <code>mediaId</code> and
	// <code>entryPoint</code> fields.
	//
	EntryPoint string `xml:"entryPoint,omitempty"`
}

type MediaBundleError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MediaBundleError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *MediaBundleErrorReason `xml:"reason,omitempty"`
}

type MediaError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MediaError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *MediaErrorReason `xml:"reason,omitempty"`
}

type MediaPage struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 MediaPage"`

	//
	// The result entries in this page.
	//
	Entries []*Media `xml:"entries,omitempty"`

	//
	// Total number of entries in the result that this page is a part of.
	//
	TotalNumEntries int32 `xml:"totalNumEntries,omitempty"`
}

type Media_Size_DimensionsMapEntry struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Media_Size_DimensionsMapEntry"`

	Key *MediaSize `xml:"key,omitempty"`

	Value *Dimensions `xml:"value,omitempty"`
}

type Media_Size_StringMapEntry struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Media_Size_StringMapEntry"`

	Key *MediaSize `xml:"key,omitempty"`

	Value string `xml:"value,omitempty"`
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

type Video struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 Video"`

	*Media

	//
	// The duration of the associated video, in milliseconds.
	// <span class="constraint Selectable">This field can be selected using the value "DurationMillis".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	DurationMillis int64 `xml:"durationMillis,omitempty"`

	//
	// Streaming URL for the video.
	// <span class="constraint Selectable">This field can be selected using the value "StreamingUrl".</span>
	//
	StreamingUrl string `xml:"streamingUrl,omitempty"`

	//
	// Indicates whether the video is ready to play on the web.
	// <span class="constraint Selectable">This field can be selected using the value "ReadyToPlayOnTheWeb".</span>
	//
	ReadyToPlayOnTheWeb bool `xml:"readyToPlayOnTheWeb,omitempty"`

	//
	// The Industry Standard Commercial Identifier code for this media, used
	// mainly for television commercials.
	// <span class="constraint Selectable">This field can be selected using the value "IndustryStandardCommercialIdentifier".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	IndustryStandardCommercialIdentifier string `xml:"industryStandardCommercialIdentifier,omitempty"`

	//
	// The Advertising Digital Identification code for this media, as defined by
	// the American Association of Advertising Agencies, used mainly for
	// television commercials.
	// <span class="constraint Selectable">This field can be selected using the value "AdvertisingId".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	AdvertisingId string `xml:"advertisingId,omitempty"`

	//
	// For YouTube-hosted videos, the YouTube video ID (as seen in YouTube URLs)
	// may also be filled in.
	// <span class="constraint Selectable">This field can be selected using the value "YouTubeVideoIdString".</span><span class="constraint Filterable">This field can be filtered on.</span>
	//
	YouTubeVideoIdString string `xml:"youTubeVideoIdString,omitempty"`
}

type VideoError struct {
	XMLName xml.Name `xml:"https://adwords.google.com/api/adwords/cm/v201802 VideoError"`

	*ApiError

	//
	// The error reason represented by an enum.
	//
	Reason *VideoErrorReason `xml:"reason,omitempty"`
}

type MediaServiceInterface struct {
	client *SOAPClient
}

func NewMediaServiceInterface(url string, tls bool, auth *BasicAuth) *MediaServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClient(url, tls, auth)

	return &MediaServiceInterface{
		client: client,
	}
}

func NewMediaServiceInterfaceWithTLSConfig(url string, tlsCfg *tls.Config, auth *BasicAuth) *MediaServiceInterface {
	if url == "" {
		url = ""
	}
	client := NewSOAPClientWithTLSConfig(url, tlsCfg, auth)

	return &MediaServiceInterface{
		client: client,
	}
}

func (service *MediaServiceInterface) AddHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Backwards-compatible function: use AddHeader instead
func (service *MediaServiceInterface) SetHeader(header interface{}) {
	service.client.AddHeader(header)
}

// Error can be either of the following types:
//
//   - ApiException
/*
   Returns a list of media that meet the criteria specified by the selector.
   <p class="note"><b>Note:</b> {@code MediaService} will not return any
   {@link ImageAd} image files.</p>

   @param serviceSelector Selects which media objects to return.
   @return A list of {@code Media} objects.
*/
func (service *MediaServiceInterface) Get(request *Get) (*GetResponse, error) {
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
   Returns the list of {@link Media} objects that match the query.

   @param query The SQL-like AWQL query string
   @returns A list of {@code Media} objects.
   @throws ApiException when the query is invalid or there are errors processing the request.
*/
func (service *MediaServiceInterface) Query(request *Query) (*QueryResponse, error) {
	response := new(QueryResponse)
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
   Uploads new media. Currently, you can upload {@link Image} files and {@link MediaBundle}s.

   @param media A list of {@code Media} objects, each containing the data to
   be uploaded.
   @return A list of uploaded media in the same order as the argument list.
*/
func (service *MediaServiceInterface) Upload(request *Upload) (*UploadResponse, error) {
	response := new(UploadResponse)
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
