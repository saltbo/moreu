package api

// 第一步：跳去登录页面，如果未登录则需要用户输入账号密码,请求后端登录接口进行登录（或跳去其他社交登录，登录成功后再跳回来；
// 第二步：登录成功后跳去oauth2/auth接口，基于登录的用户信息生成Session并302一个code给浏览器端
// 第三步：浏览器端拿着code去请求业务callback接口
// 第四步：业务callback接口拿着code向oauth2/token接口换取access_token和id_token

