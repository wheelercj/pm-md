# calendar API v1

----------------------------------------

## get all accounts

GET `/v1/account/all`

### sample response (status: 200 OK)

```json
[
    {
        "_id": "64de85af99bb7d63123531e8",
        "email": "alksdfuie@mail.com",
        "hashedPassword": "$2b$10$1Pkuaf10UTOXaY8WctU72em4HDOHiAVLxssXc2iqIEz0BbWEE/g5q",
        "scheduledAppointmentCount": 0,
        "editedAppointmentCount": 0,
        "canceledAppointmentCount": 0,
        "createdAt": "2023-08-17T20:40:15.775Z",
        "activeAppointments": [],
        "__v": 0
    }
]
```

----------------------------------------

## create account

POST `/v1/account/register`

### sample request body

```json
{
    "email": "alksdfuie@mail.com",
    "password": "so8vu3os8srl3jmlsf"
}
```

### sample response (status: 201 Created)

```json
{
    "email": "alksdfuie@mail.com"
}
```

----------------------------------------

## log in

POST `/v1/account/login`

### sample request body

```json
{
    "email": "alksdfuie@mail.com",
    "password": "so8vu3os8srl3jmlsf"
}
```

### sample response (status: 200 OK)

```json
{
    "email": "alksdfuie@mail.com"
}
```

----------------------------------------

## edit account

PUT `/v1/account`

### sample request body

```json
{
    "email": "alksdfuie@mail.com",
    "password": "so8vu3os8srl3jmlsf",
    "newEmail": "alksdfuie@mail.com",
    "newPassword": "bfls83uxlf3lajbla"
}
```

### sample response (status: 200 OK)

```html
Update successful
```

----------------------------------------

## delete account

DELETE `/v1/account`

### sample request body

```json
{
    "email": "alksdfuie@mail.com",
    "password": "bfls83uxlf3lajbla"
}
```

### sample response (status: 200 OK)

```html
Account deleted
```

