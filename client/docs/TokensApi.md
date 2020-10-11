# {{classname}}

All URIs are relative to *//localhost:8080/moreu/api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**TokensDelete**](TokensApi.md#TokensDelete) | **Delete** /tokens | 退出登录
[**TokensPost**](TokensApi.md#TokensPost) | **Post** /tokens | 登录/密码重置

# **TokensDelete**
> HttputilJsonResponse TokensDelete(ctx, )
退出登录

用户状态登出

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

# **TokensPost**
> HttputilJsonResponse TokensPost(ctx, body)
登录/密码重置

用于账户登录和申请密码重置

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**BindBodyToken**](BindBodyToken.md)| 参数 | 

### Return type

[**HttputilJsonResponse**](httputil.JSONResponse.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

