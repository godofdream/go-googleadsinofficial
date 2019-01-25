#!/bin/bash

if [ ! -x gowsdl ]; then
    go get -u github.com/hooklift/gowsdl/...
fi

export PATH="$PATH:$1/bin"
export version="v201809"
mkdir -p $version
cd $version
for url in \
https://adwords.google.com/api/adwords/mcm/$version/AccountLabelService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/AdCustomizerFeedService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/AdGroupAdService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/AdGroupBidModifierService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/AdGroupCriterionService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/AdGroupExtensionSettingService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/AdGroupFeedService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/AdGroupService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/AdParamService?wsdl \
https://adwords.google.com/api/adwords/rm/$version/AdwordsUserListService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/BatchJobService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/BiddingStrategyService?wsdl \
https://adwords.google.com/api/adwords/billing/$version/BudgetOrderService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/BudgetService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/CampaignBidModifierService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/CampaignCriterionService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/CampaignExtensionSettingService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/CampaignFeedService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/CampaignGroupPerformanceTargetService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/CampaignGroupService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/CampaignService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/CampaignSharedSetService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/ConstantDataService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/ConversionTrackerService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/CustomerExtensionSettingService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/CustomerFeedService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/CustomerNegativeCriterionService?wsdl \
https://adwords.google.com/api/adwords/mcm/$version/CustomerService?wsdl \
https://adwords.google.com/api/adwords/ch/$version/CustomerSyncService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/DataService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/DraftAsyncErrorService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/DraftService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/FeedItemService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/FeedItemTargetService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/FeedMappingService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/FeedService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/LabelService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/LocationCriterionService?wsdl \
https://adwords.google.com/api/adwords/mcm/$version/ManagedCustomerService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/MediaService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/OfflineCallConversionFeedService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/OfflineConversionFeedService?wsdl \
https://adwords.google.com/api/adwords/rm/$version/OfflineDataUploadService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/ReportDefinitionService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/SharedCriterionService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/SharedSetService?wsdl \
https://adwords.google.com/api/adwords/o/$version/TargetingIdeaService?wsdl \
https://adwords.google.com/api/adwords/o/$version/TrafficEstimatorService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/TrialAsyncErrorService?wsdl \
https://adwords.google.com/api/adwords/cm/$version/TrialService?wsdl \
; do
  service1="${url##*/}"
  service="${service1:0:(-5)}"
  wget -O ${service}.wsdl "$url"

  #mkdir -p ${service}
  #wsdl2go -p googleadsinofficial -i ${service}.wsdl -o ${service}/${service}.go
  gowsdl -o ${service}.go -p ${service} ${service}.wsdl
done
