# Changelog

## [1.1.1](https://github.com/Martins-Iroka/MyGallery-Backend/compare/MyGallery-Backend-v1.1.0...MyGallery-Backend-v1.1.1) (2025-11-08)


### Bug Fixes

* remove example logging of photographer names in main.go and update release-type in release-please config to 'go' ([7479fd8](https://github.com/Martins-Iroka/MyGallery-Backend/commit/7479fd8d2016a4f516b18dcd249d4b32ada4341d))

## [1.1.0](https://github.com/Martins-Iroka/MyGallery-Backend/compare/MyGallery-Backend-v1.0.0...MyGallery-Backend-v1.1.0) (2025-11-08)


### Features

* add build target to Makefile ([9cd446b](https://github.com/Martins-Iroka/MyGallery-Backend/commit/9cd446bcfbc12dcb4b4124c435629dcc7ce83c98))
* add CI/CD workflows for build, release, and version update automation ([cf204aa](https://github.com/Martins-Iroka/MyGallery-Backend/commit/cf204aafec07c6490056cdbb555c2a8693962f1b))
* add initial structures for photo and video storage, including Pexels response handling ([242dc82](https://github.com/Martins-Iroka/MyGallery-Backend/commit/242dc82e0440bb03ed4994aba823c3a587d02f77))
* add manifest file for release-please and update workflow configuration ([dbcc115](https://github.com/Martins-Iroka/MyGallery-Backend/commit/dbcc1153fd6c229572f1e7363f675215c6f8103a))
* add migration scripts for photo posts, video posts, and video download files tables ([88b844f](https://github.com/Martins-Iroka/MyGallery-Backend/commit/88b844f9938b9c2ae918ce3a6c5dab7dc22ed14b))
* add user login endpoint and request payload validation ([937560f](https://github.com/Martins-Iroka/MyGallery-Backend/commit/937560f78dc2724309e335b38965455df5780e61))
* add user verification endpoint and payload structures ([d36fdc3](https://github.com/Martins-Iroka/MyGallery-Backend/commit/d36fdc38fc3057d79fb8452f3590d90a7b940ceb))
* added air.toml ([a64791c](https://github.com/Martins-Iroka/MyGallery-Backend/commit/a64791c5058e729ef1a4b1c88f887af27dc66302))
* added docker-compose.yml for db setup ([4d0a612](https://github.com/Martins-Iroka/MyGallery-Backend/commit/4d0a612bfa137fab9f576c1dbdb7a97f01444346))
* added README ([2fe5f99](https://github.com/Martins-Iroka/MyGallery-Backend/commit/2fe5f99c667565e33396864a8d95e216abc0c386))
* added userstore to communicate with db ([298f64f](https://github.com/Martins-Iroka/MyGallery-Backend/commit/298f64f0e764b49a92773566f47ad88a872cbe96))
* implement API call to fetch curated photos from Pexels ([dbd256f](https://github.com/Martins-Iroka/MyGallery-Backend/commit/dbd256f0f6cdcfbe740794118e3540b9ead3db6d))
* implement asynchronous email verification with SAGA compensation in user registration ([2193140](https://github.com/Martins-Iroka/MyGallery-Backend/commit/2193140e54c2359ed387da5658f82980e959047b))
* implement main function in seed migration ([773443e](https://github.com/Martins-Iroka/MyGallery-Backend/commit/773443e4634f1451a8b5199edbb0b9ac6ca33c8f))
* implement OTP verification and refactor user registration endpoints, enhance validation and testing ([2ddab78](https://github.com/Martins-Iroka/MyGallery-Backend/commit/2ddab78459d70246ae8a81a74b830004c1e2cbea))
* implement user login functionality with JWT authentication and password hashing ([870c23b](https://github.com/Martins-Iroka/MyGallery-Backend/commit/870c23b13f5d09fb5f20adea05c674dda178b6dd))
* implemented db connection logic and confirmed connection ([da38e30](https://github.com/Martins-Iroka/MyGallery-Backend/commit/da38e3009da52750d64099ca7180530cd2d944d5))
* implemented jwt authenticator and test ([d7e990d](https://github.com/Martins-Iroka/MyGallery-Backend/commit/d7e990dd4c35a4ca0cba663ab2c7f9dc3c67907a))
* refactor error handling and user storage logic, migrate utility functions to a new package ([7ed7ee3](https://github.com/Martins-Iroka/MyGallery-Backend/commit/7ed7ee30420ade051f9e33bf5a40b0c18cb6cb61))
* refactor user creation and verification logic, update method names and implement transaction handling ([204bae0](https://github.com/Martins-Iroka/MyGallery-Backend/commit/204bae04d981c8dc0e7485854c2cdbeebd0454c4))
* restructure application to support user registration and verification via Twilio ([6342746](https://github.com/Martins-Iroka/MyGallery-Backend/commit/63427467985a08ea21fb7ddd2389b2b5eb5e67c1))
* setup app config and router ([f0e9be7](https://github.com/Martins-Iroka/MyGallery-Backend/commit/f0e9be7e2048941c2f0a18df65811c0497227d11))
* update user registration and verification payloads, enhance Swagger documentation, and modify build configuration ([6b111fd](https://github.com/Martins-Iroka/MyGallery-Backend/commit/6b111fd94391c5d0880382dbd3efa0a55f884a30))


### Bug Fixes

* add missing newline at end of file in release-please.yaml ([7e0a68a](https://github.com/Martins-Iroka/MyGallery-Backend/commit/7e0a68ade71540890ab2398e7a1014961953254e))
* correct release-type key format in release-please config and remove redundant line in workflow ([403f997](https://github.com/Martins-Iroka/MyGallery-Backend/commit/403f997802034c762ddca245d91bb576045ab106))
* revert release-type in release-please config from 'go' to 'simple' ([d817638](https://github.com/Martins-Iroka/MyGallery-Backend/commit/d817638193d759f9ee5a48c63f47e77fe80b4da7))
* suppress unused error variable in JWT tests ([85c9791](https://github.com/Martins-Iroka/MyGallery-Backend/commit/85c97915dbcf91731fda467294424308bcc085f6))
* update config file path in release-please workflow ([7090ae6](https://github.com/Martins-Iroka/MyGallery-Backend/commit/7090ae62245ac21e1201919cb9044ec9e259e729))
* update release-please configuration by removing unnecessary fields and restructuring manifest ([75d678a](https://github.com/Martins-Iroka/MyGallery-Backend/commit/75d678a546b382561c3318a3015f27dd175c1357))
* update release-please configuration to standard key format and ensure proper JSON structure ([2b4d6d5](https://github.com/Martins-Iroka/MyGallery-Backend/commit/2b4d6d531243e4001d37d7fb7941856eeb30d073))
* update release-type in release-please manifest to 'go' ([2da9443](https://github.com/Martins-Iroka/MyGallery-Backend/commit/2da94437705f99a8cff043cd8b3244b62d802fbf))
