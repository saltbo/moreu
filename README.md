# moreu

[![](https://github.com/saltbo/moreu/workflows/build/badge.svg)](https://github.com/saltbo/moreu/actions?query=workflow%3Abuild)
[![](https://codecov.io/gh/saltbo/moreu/branch/master/graph/badge.svg)](https://codecov.io/gh/saltbo/moreu)
[![](https://wakatime.com/badge/github/saltbo/moreu.svg)](https://wakatime.com/badge/github/saltbo/moreu)
[![](https://api.codacy.com/project/badge/Grade/88817db9b3b04c0293c9d001d574a5ef)](https://app.codacy.com/manual/saltbo/moreu?utm_source=github.com&utm_medium=referral&utm_content=saltbo/moreu&utm_campaign=Badge_Grade_Dashboard)
[![](https://img.shields.io/github/v/release/saltbo/moreu.svg)](https://github.com/saltbo/moreu/releases)
[![](https://img.shields.io/github/license/saltbo/moreu.svg)](https://github.com/saltbo/moreu/blob/master/LICENSE)

English | [🇨🇳中文](https://saltbo.cn/moreu)

## Endpoints
<!-- markdown-swagger -->
 Endpoint            | Method | Auth? | Description               
 ------------------- | ------ | ----- | --------------------------
 `/tokens`           | POST   | No    | 用于账户登录和申请密码重置
 `/users`            | GET    | No    | 获取一个用户信息          
 `/users`            | POST   | No    | 注册一个用户              
 `/users/{email}`    | PATCH  | No    | 用于账户激活和密码重置    
 `/users/{username}` | GET    | No    | 获取一个用户信息          
<!-- /markdown-swagger -->

## Contributing
See [CONTRIBUTING](CONTRIBUTING.md) for details on submitting patches and the contribution workflow.

## Contact us
- [Author Blog](https://saltbo.cn).

## Author
- [saltbo](https://github.com/saltbo)

## License
- [MIT](https://github.com/saltbo/moreu/blob/master/LICENSE)