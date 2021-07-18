# appointment

Microservice to be used by doctors and patients to manage appointments

Appointment lets Doctor create schedule for their availability. Patient can check Doctor's schedule and book an Appointment if slot is available. Doctor can only create schedule for the current day. The size of a slot is 15 mins.

### Endpoints

/schedule : Doctor can use this to specify what time he/she is available for appointments.

/book : Used by Patient to book the 15 min time slot with the Doctor.

/list : Used to list the complete schedule of the Doctor's appointments for the day.

/cancel : Used to cancel an appointment. Can be used by either Doctor or Patient.

/signup : Used to signup for the service and recieve a token which will be required in all further interactions.

<br/> <br/>
**N.B**
Listening port of the service can be configured by using the **PORT** environment variable. defaults to 8080.

<br/> <br/>

## Usage

All endpoints accept valid JSON and respond with valid JSON.

All Time values to be provided in "YYYY-mm-ddTHH:MM:SSZ" format only. e.g. 2021-07-18T13:30:00Z
<br/> <br/>

### /signup

---

The first step is to hit the /signup endpoint to get a Token

#### Input Format:

```json
{
  "name": "Sachin",
  "usertype": "Doctor"
}
```

#### Fields:

- **name (String)** : Name of the User

- **usertype (String)** : Type of User. Allowed values - "Patient" or "Doctor"

#### Output Format:

```json
{
  "message": "Account created",
  "status": 200,
  "token": "MXxEb2N0b3I"
}
```

<br/>

### /schedule

---

Next the Doctor can create his/her schedule

#### Input Format:

```json
{
  "starttime": "2021-07-18T19:00:00Z",
  "endtime": "2021-07-18T20:30:00Z",
  "token": "MXxEb2N0b3I"
}
```

#### Fields:

- **starttime (Time)** : Start time of schedule

- **endtime (Time)** : End time of schedule

- **token** : Token generated in Step 1

#### Output Format:

```json
{
  "message": "Schedule created",
  "status": 200
}
```

<br/>

### /list

---

Patient can list Doctor's schedule

#### Input Format:

```json
{
  "doctorname": "Sachin",
  "token": "MXxQYXRpZW50"
}
```

#### Fields:

- **doctorname (String)** : Name of the doctor to list schedule for

- **token** : Token generated in Step 1

#### Output Format:

```json
{
  "appointments": [
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:00:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:15:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:30:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:45:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T20:00:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T20:15:00Z",
      "booked": false
    }
  ],
  "message": "Appointments Listed",
  "status": 200
}
```

<br/>

### /book

---

Patient can book an Appointment

#### Input Format:

```json
{
  "doctorname": "Sachin",
  "starttime": "2021-07-18T19:30:00Z",
  "token": "MXxQYXRpZW50"
}
```

#### Fields:

- **doctorname (String)** : Name of the doctor whose appointment to be booked

- **starttime (Time)** : Start time of appointment

- **token** : Token generated in Step 1

#### Output Format:

```json
{
  "appointmentid": 1,
  "message": "Appointment booked",
  "status": 200
}
```

- **appointmentid** : To be used while cancelling appointment

<br/>
Another Patient cant book the same slot

##### Input:

```json
{
  "doctorname": "Sachin",
  "starttime": "2021-07-18T19:30:00Z",
  "token": "MnxQYXRpZW50"
}
```

#### Output:

```json
{
  "message": "Slot already taken",
  "status": 400,
  "error": ""
}
```

<br/>
They can book a different slot

##### Input:

```json
{
  "doctorname": "Sachin",
  "starttime": "2021-07-18T20:00:00Z",
  "token": "MnxQYXRpZW50"
}
```

#### Output:

```json
{
  "appointmentid": 2,
  "message": "Appointment booked",
  "status": 200
}
```

<br/>
Schedule listing after 2 slots were booked

```json
{
  "appointments": [
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:00:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:15:00Z",
      "booked": false
    },
    {
      "appointmentid": "1",
      "doctorid": "1",
      "patientid": "1",
      "starttime": "2021-07-18T19:30:00Z",
      "booked": true
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:45:00Z",
      "booked": false
    },
    {
      "appointmentid": "2",
      "doctorid": "1",
      "patientid": "2",
      "starttime": "2021-07-18T20:00:00Z",
      "booked": true
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T20:15:00Z",
      "booked": false
    }
  ],
  "message": "Appointments Listed",
  "status": 200
}
```

<br/>

### /cancel

---

Patient can cancel their appointment

#### Input Format:

```json
{
  "appointmentid": 1,
  "token": "MXxQYXRpZW50"
}
```

#### Fields:

- **appointmentid (Int)** : Appointment ID received when slot was booked

- **token** : Token generated in Step 1

#### Output Format:

```json
{
  "message": "Appointment cancelled",
  "status": 200
}
```

<br/>
Schedule listing after 1 cancellation

```json
{
  "appointments": [
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:00:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:15:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:30:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:45:00Z",
      "booked": false
    },
    {
      "appointmentid": "2",
      "doctorid": "1",
      "patientid": "2",
      "starttime": "2021-07-18T20:00:00Z",
      "booked": true
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T20:15:00Z",
      "booked": false
    }
  ],
  "message": "Appointments Listed",
  "status": 200
}
```

<br/>
Doctor can also cancel an appointment

#### Input:

```json
{
  "appointmentid": 2,
  "token": "MXxEb2N0b3I"
}
```

#### Output:

```json
{
  "message": "Appointment cancelled",
  "status": 200
}
```

<br/>
Schedule listing after both cancellations

```json
{
  "appointments": [
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:00:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:15:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:30:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T19:45:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T20:00:00Z",
      "booked": false
    },
    {
      "appointmentid": "",
      "doctorid": "1",
      "patientid": "",
      "starttime": "2021-07-18T20:15:00Z",
      "booked": false
    }
  ],
  "message": "Appointments Listed",
  "status": 200
}
```
