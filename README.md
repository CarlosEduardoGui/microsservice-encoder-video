# Video Encoder Microservice

## Environment Setup

To run in development mode, follow these steps:

* Duplicate the file `.env.example` to `.env`
* Run the `docker-compose up -d` command
* Access the RabbitMQ administration and create an exchange of type `fanout`. It will serve as a `Dead Letter Exchange` to receive messages that are not processed.
* Create a `Dead Letter Queue` and bind it to the `Dead Letter Exchange` that was just created. There is no need for a routing_key.
* In the `.env` file, specify the name of the `Dead Letter Exchange` in the parameter: `RABBITMQ_DLX`
* Create a service account on GCP that has permission to write to Google Cloud Storage. Download the JSON file with the credentials and save it in the project's root exactly with the name: `bucket-credential.json`


## Executing

To execute the encoder, run the `make server` command directly in the container. Example:

```
docker exec encoder-new2_app_1 make server
```

Where `encoder-new2_app_1` is the name of the container generated by docker-compose.

## Message Sending Pattern to the Encoder

For a message to be parsed by the encoder system, it must arrive in the following JSON format:

```
{
  "resource_id": "my-resource-id-can-be-a-uuid-type",
  "file_path": "convite.mp4"
}
```

* `resource_id`: Represents the ID of the video you want to convert. It is of string type.
* `file_path`: Is the full path of the mp4 video within the bucket.

## Message Return Pattern by the Encoder

### Success in processing

For each processed video, the encoder will send the processing result to an exchange (to be configured in .env).

If the processing has been completed successfully, the JSON return pattern will be:

```
{
    "id":"bbbdd123-ad05-4dc8-a74c-d63a0a2423d5",
    "output_bucket_path":"codeeducationtest",
    "status":"COMPLETED",
    "video":{
        "encoded_video_folder":"b3f2d41e-2c0a-4830-bd65-68227e97764f",
        "resource_id":"aadc5ff9-0b0d-13ab-4a40-a11b2eaa148c",
        "file_path":"convite.mp4"
    },
    "Error":"",
    "created_at":"2020-05-27T19:43:34.850479-04:00",
    "updated_at":"2020-05-27T19:43:38.081754-04:00"
}
```

Being that `encoded_video_folder` is the folder that contains the converted video.

### Error in processing

If the processing has encountered any error, the JSON return pattern will be:

```
{
    "message": {
        "resource_id": "aadc5ff9-010d-a3ab-4a40-a11b2eaa148c",
        "file_path": "invite.mp4"
    },
    "error":"Reason for the error"
}
```

In addition, the encoder will send the original message that had a problem during processing to a dead letter exchange.
Just configure the desired DLX in the .env file in the parameter: `RABBITMQ_DLX`

## Complete application diagram.

![PlantUML's solution](https://www.planttext.com/api/plantuml/svg/jLVBRkCs5DtxArXX5WwWn9kk0XGOnybqqYQDOtlQfN0jqJ9hYHH8AXvFqN_l9IdPrClE5EqiCXP5tdDuxbCVbGQfopmZiU2QVjLYz0FZSLOcb6orBjJjP29XVvColItNfQBIxkFw9XRvfCY0cyFNIYSPMqzcFrxFpTANvwT93afJYKlH34y0urJU5BXtD4sI1SrMa8u3K0SU1o7vaE7hfJvabFn4xa5jQVA4A-Ehgz7oyyiVqVnXiNpvBPg_B5QBwUDgV3KQQQPnTp0J2TraxEijcKk8DOiXmr_YJnfh8ZP4RodvL6OPcHrYt3uJfLGq6CsfL8AFOz2GYJI296USAIr12na6CxTCVcS95MJglmp2u1auf5IHD_DM4SOjI6MA04_CCvDcXQP2Cgg2gnfGbMiLIxQciesvFBuP15G753Gd9nig9A38dkBryFzglT1CS68SquiR1kWWlH4o9oV8mCDNH0jxewBUXyyGnrJLLqQpewNb4js9aYSW3-CKBxk7GtXoI715_XeH_7Vq_dWZRUWaKKFtCFXKhTzMYjVwMo6GsuAYPmt7-HIPWIPvmmKXSG7BUEg9RZsHKYwj_aqYQuXtFJpktF8GZIVWQOcurv1oeJ3KFDyOmIoREX3-8r9HWazhK3GCxgjVZey9xWPDG-QmV7qoEc7S89OvMkT0Wri5ZKauv981WHvnhyao33ei1sWK4EtpQZLVJdx1_FJizXkuTCMM_LAXUeYmciqhqot9A8ynsZ2cdH8FxX5F_2pm3yQ-YTygJ8ZTRDo1YoGlQ7gTOA7fkg54ZZggGYkuiu_sZ49PpNc7cjJPW6X7VWDc8goWlQPmlS-CKR2-MulFbo8Ja2XJhT-er3KDtnUNZ3WszsBFdQaZYG4Z0xxicgov3HdiqkeFCrHHbRaNmPpkimE99ONsXHUmNZFz_4fUIrfMOQvkS7yGVj2GAcmdLPFfYa86uEVbEcFOWfGhDWvk3qtDXy6FkyvGw4ScntZkAD3lxl-FG4i3C2xmb50u0DHzmdZR1SE_yCc53b5dTbsAwqxXDGLizWP14Gvl3DkGwNwcoZgosCGEvgukwqwTivElKB0XqEdySQ3w7x6HKY4Ony9bSCuowI3J38mvHK76swZ1esPmAlsMoeD_ebllwwxe-ao8LLzXowHvDJ4WeyPeqjO1uYsDTh3PGHXCxHFvrJoP1eD304f0xgJOkQd18HkRE1c-rIHmDBI5LehSdJK074Uo3fp5YMBtDWoQcrtF8wOhfDBgQI0HXyNK2YrfaU74fwicQxKo1Bn0z_Kye_JxErkJyHcJImmArMoyJXrdVCduZhGyEPYkJT1BRNGc2bQxmlPGKfsqu-GzNTDlhCn9DIwy4xGP6BrwEUzXKcWHYsmesriHu8-j_Lcsz7Pu8vSTSLteUHkKoO1pYUvlZWh4I4sID42E-MPIlDhHJJCFn9lwwy6x5h5nlG7iNZ5hu9EUdMgve_XH8s27obU_yDnZWhMPwGEf4iKVXA14abpRbC5hcXG2_7WuJlhcyqxvezQfPVT_AWrDzCFfv5ltz35V9ir8k6YYTVeanlUST63lF7jfQZwPcPw4fkOJiS-8nxtRzjvrXywnxonfFLHQ_Ztd4e_GtwEFIAFCity0)
