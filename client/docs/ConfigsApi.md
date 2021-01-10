# {{classname}}

All URIs are relative to *//localhost:8081/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ConfigsGet**](ConfigsApi.md#ConfigsGet) | **Get** /configs | 获取全部配置
[**ConfigsKeyDelete**](ConfigsApi.md#ConfigsKeyDelete) | **Delete** /configs/{key} | 删除配置项
[**ConfigsKeyGet**](ConfigsApi.md#ConfigsKeyGet) | **Get** /configs/{key} | 获取配置项
[**ConfigsPost**](ConfigsApi.md#ConfigsPost) | **Post** /configs | 创建配置项

# **ConfigsGet**
> HttputilJsonResponse ConfigsGet(ctx, )
获取全部配置

获取全部配置

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**HttputilJsonResponse**](httputil.JSONResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

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

# **ConfigsPost**
> HttputilJsonResponse ConfigsPost(ctx, body)
创建配置项

创建配置项

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**BindBodyConfig**](BindBodyConfig.md)| 参数 | 

### Return type

[**HttputilJsonResponse**](httputil.JSONResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

