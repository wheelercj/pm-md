# calendar API

----------------------------------------

<details open>
<summary>
<h1>POST endpoints</h1>
</summary>

----------------------------------------

<details open>
<summary>
<h2>create account</h2>
</summary>

POST `/v1/account/register`

<h3>sample request body</h3>

```json
{
    "email": "alksdfuie@mail.com",
    "password": "so8vu3os8srl3jmlsf"
}
```

<details>
    <summary>
        <h3>sample response to valid input (status: 201 Created)</h3>
    </summary>

```json
{
    "email": "alksdfuie@mail.com"
}
```
</details>
</details>

----------------------------------------

<details open>
<summary>
<h2>log in</h2>
</summary>

POST `/v1/account/login`

<h3>sample request body</h3>

```json
{
    "email": "alksdfuie@mail.com",
    "password": "so8vu3os8srl3jmlsf"
}
```

<details>
    <summary>
        <h3>sample response to valid input (status: 200 OK)</h3>
    </summary>

```json
{
    "email": "alksdfuie@mail.com"
}
```
</details>
</details>
</details>

----------------------------------------

<details open>
<summary>
<h1>empty folder</h1>
</summary>
</details>

----------------------------------------

<details open>
<summary>
<h1>GET endpoints</h1>
</summary>

----------------------------------------

<details open>
<summary>
<h2>get all accounts</h2>
</summary>

GET `/v1/account/all`

<details>
    <summary>
        <h3>sample response to valid input (status: 200 OK)</h3>
    </summary>

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
</details>
</details>
</details>

----------------------------------------

<details open>
<summary>
<h1>edit account</h1>
</summary>

PUT `/v1/account`

<h2>sample request body</h2>

```json
{
    "email": "alksdfuie@mail.com",
    "password": "so8vu3os8srl3jmlsf",
    "newEmail": "alksdfuie@mail.com",
    "newPassword": "bfls83uxlf3lajbla"
}
```

<details>
    <summary>
        <h2>sample response to valid input (status: 200 OK)</h2>
    </summary>

```html
Update successful
```
</details>
</details>

----------------------------------------

<details open>
<summary>
<h1>delete account</h1>
</summary>

DELETE `/v1/account`

<h2>sample request body</h2>

```json
{
    "email": "alksdfuie@mail.com",
    "password": "bfls83uxlf3lajbla"
}
```

<details>
    <summary>
        <h2>sample response to valid input (status: 200 OK)</h2>
    </summary>

```html
Account deleted
```
</details>
</details>
