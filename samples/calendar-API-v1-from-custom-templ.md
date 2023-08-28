# calendar API v1

----------------------------------------

## get all accounts

GET `/v1/account/all`

### sample response "get all accounts" (status: 200 OK)

```json
[
    {
        "_id": "64d95ea1a74725d232891e29",
        "email": "alksdfuie@mail.com",
        "hashedPassword": "$2b$10$l/DUytcwyH5fgw7uTbl0OuGbXcFst4Yqdu7Ueh73GhdQYBnkAkiCC",
        "scheduledAppointmentCount": 0,
        "editedAppointmentCount": 0,
        "canceledAppointmentCount": 0,
        "createdAt": "2023-08-13T22:52:17.292Z",
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

### sample response "create account" (status: 201 Created)

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

### sample response "log in" (status: 200 OK)

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

### sample response "edit account" (status: 200 OK)

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

### sample response "delete account" (status: 200 OK)

```html
Account deleted
```

