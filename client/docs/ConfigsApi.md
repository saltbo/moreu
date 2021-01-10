# {{classname}}

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ConfigsKeyDelete**](ConfigsApi.md#ConfigsKeyDelete) | **Delete** /configs/{key} | 删除配置项
[**ConfigsKeyGet**](ConfigsApi.md#ConfigsKeyGet) | **Get** /configs/{key} | 获取配置项
[**ConfigsKeyPut**](ConfigsApi.md#ConfigsKeyPut) | **Put** /configs/{key} | 修改配置项

# **ConfigsKeyDelete**
> HttputilJsonResponse ConfigsKeyDelete(ctx, key)
删除配置项

根据键名删除配置项

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **key** | **string**| 键名 | 

### Return type

[**HttputilJsonResponse**](httputil.JSONResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ConfigsKeyGet**
> HttputilJsonResponse ConfigsKeyGet(ctx, key)
获取配置项

根据键名获取配置项

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **key** | **string**| 键名 | 

### Return type

[**HttputilJsonResponse**](httputil.JSONResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ConfigsKeyPut**
> HttputilJsonResponse ConfigsKeyPut(ctx, body, key)
修改配置项

根据键名修改配置项

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**interface{}**](interface{}.md)| 参数 | 
  **key** | **string**| 键名 | 

### Return type

[**HttputilJsonResponse**](httputil.JSONResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

