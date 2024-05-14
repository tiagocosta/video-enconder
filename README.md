# video-enconder
This is a video encoder service that downloads a video file from a bucket, fragments and packages (encodes) it and uploads the encoded files back to the bucket.

Basically, the service consumes a queue and waits for a video request from any other service/client. When a video request arrives, the service communicates with the cloud bucket (GCP - Cloud Storage) to download the requested video from there. 
Once the video file is downloaded it starts the process of fragmentation and packaging/encoding. After this process is completed many encoded files are uploaded to the bucket.
Finally, a message is published in another queue (or in a dead letter queue in case of fail) to make the requesting service know about the result.

## Dependencies
Go 1.22.2

Bento4 1.6.0-637

MySQL 5.7

RabbitMQ
