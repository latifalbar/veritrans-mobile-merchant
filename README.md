# Reference Merchant Server for Veritrans Mobile SDK

This is a testing server for the development of Veritrans Mobile SDK. Also acts as a reference implementation for the methods to be implemented by merchants to use the mobile SDK.

## Implementation

Testing server will provide two different API endpoint for `Sandbox` and `Production` build.

- `Sandbox` will use `/api` endpoint group.
- `Production` will use `/api-prod` endpoint group.

Each endpoint group will have several function which needed to run by Veritrans Mobile SDK.

### Charge Payment

- **Endpoint**: `POST /api/charge`
- **Headers**:
    - Content-Type  : `application/json`
    - Accept        : `application/json`
    - Authorization : `Basic ` + Base64(__VT-Server-Key__)
- **Request**: All Veritrans supported payment details. Read more at [Veritrans Documentation Website](http://docs.veritrans.co.id/en/api/methods.html#Charge).
- **Response**: Should be same as the response from Veritrans Payment API. Read more at [Veritrans Documentation Website](http://docs.veritrans.co.id/en/api/methods.html#Charge).

This endpoint will be used to charge payment to Veritrans Payment API. Basically the implementation of this endpoint just add `Authorization` in the header using `Server Key` which is required to do the charging and then do a HTTP call to Veritrans Payment API. The response of the HTTP call should be passed back as the response of this endpoint.

### Promotions

To handle promotion, there are several endpoint to be implemented.

#### Add a Discount Promo

- **Endpoint**: `POST /api/promotions/discount`
- **Headers**
    - Admin-Token: __MERCHANT_ADMIN_TOKEN__
- **Request**

    ```
    {
        "title": "Installment Title",
        "description": "Installment description",
        "discount_percentage": 35,
        "bins": [
            "52111"
        ],
        "installment_terms": [
            "6"
        ]
    }
    ```

- **Response**

    ```
    {
      "status_code" : 201,
      "status_message" : "Discount created."
    }
    ```

#### Remove a Discount Promo

**To be written**

#### Add an Installment Promo

- **Endpoint**: `POST /api/promotions/installment`
- **Headers**
    - Admin-Token: __MERCHANT_ADMIN_TOKEN__
- **Request**

    ```
    {
        "title": "Discount Title",
        "description": "Discount description",
        "discount_percentage": 35,
        "bins": [
            "52111"
        ]
    }
    ```

- **Response**

    ```
    {
      "status_code" : 201,
      "status_message" : "Installment created."
    }
    ```

#### Remove an Installment Promo

**To be written**

#### Get Promotions List

- **Endpoint** : `GET /api/promotions`
- **Response** :

    ```
    {
        "data": {
            "installment": [
                INSTALLMENT_OBJECT,
                INSTALLMENT_OBJECT
            ],
            "discount": [
                DISCOUNT_OBJECT,
                DISCOUNT_OBJECT
            ]
        },
        "status_code": 200,
        "status_message": "success"
    }
    ```

### Save Card

Mobile SDK enables the user to save their card credential in a token that can be saved so it can be used later.

So merchant server needs to provide endpoint to Save card, get card list and remove card.

#### Save Card

- **Endpoint** : `POST /api/card/register`
- **Request**

    ```
    {
        "status_code" : "200",
        "masked_card" : "481111-14",
        "saved_token_id" : "4ioasfaslk490asfakj"
    }
    ```
    
- **Response**

    ```
    {
        "status_code" : 201,
        "status_message" : "Card is saved."
    }
    ```

#### Get Card List

- **Endpoint** : "GET /api/card"
- **Headers**
    - X-Auth : __Authentication-Token__
- **Response**

    ```
    {
      "data": [
        {
          "saved_token_id": "4ioasfaslk490asfakj",
          "masked_card": "481111-14"
        }
      ],
      "status_code": 200,
      "status_message": "Success"
    }
    ```

#### Remove Card

**To be written**

## Credentials

You need to use your own Veritrans credential provided in MAP.

In `merchant.go` you need to change these variables.

```
// Sandbox Token
var Token = "VT-server-MFoCwh-MkpoSlOMdqiWCx9WB"

// Production Token
var ProductionToken = "VT-server-nt-LofCmn8UhgUcK8T8Za2ep"
```

## Build and Deployment

This implementation is using Go with specific build mode into Google AppEngine.

### Requirement

These are required to build this application.

- Go 1.6+
- Google AppEngine Go SDK
- Google Cloud SDK (for deployment)

### How to Build

To build and test this application, you can use Google AppEngine GoSDK.

Go to your application directory in terminal, then do this.
```
> goapp serve
```

### Deployment to AppEngine

If you're ready to use this app on production you can deploy it on your Google AppEngine project.

You can use Google Cloud SDK to deploy this application to your Google AppEngine project.

```
> appcfg.py -A GOOGLE_PROJECT_ID -V v1 update ./
```

## License

```
Copyright 2016 Raka Westu Mogandhi

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
```
