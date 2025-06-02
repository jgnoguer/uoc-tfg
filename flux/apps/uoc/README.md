# Secrets

kubectl create secret generic mqtt-credentials --from-file=mqtt-credentials.properties

kubectl create secret generic telegram-credentials --from-file=telegram-credentials.properties

mosquitto_pub -m "{ "msg": "message from mosquitto_pub client"} " -t "sensor-topic" -u jgnoguer -h 10.43.26.129 -P uocAn1m4ls


# Devices

v4l2-ctl --list-devices

sudo apt install v4l-utils
v4l2-ctl --list-devices

v4l2-ctl --list-formats

   12  v4l2-ctl --list-devices
   13  ffmpeg -hide_banner -f video4linux2 -list_formats all -i /dev/video0
   14  apt install fmpeg
   15  apt install ffmpeg
   16  ffmpeg -hide_banner -f video4linux2 -list_formats all -i /dev/video0
   17  v4l2-ctl --list-devices
   18  ffmpeg -hide_banner -f video4linux2 -list_formats all -i /dev/media0
   19  ffmpeg -hide_banner -f video4linux2 -list_formats all -i /dev/video0
   20  ffmpeg -f v4l2 -s 1280x720 -i /dev/video0  -ss 1 -frames 1 webcam-image-capture.png
   21  v4l2-ctl --list-devices
   22  ffmpeg -f v4l2 -s 1280x720 -i /dev/video0  -ss 1 -frames 1 webcam-image-capture.png
   23  v4l2-ctl --list-devices
   24  ffmpeg -f v4l2 -s 1280x720 -i /dev/video0  -ss 1 -frames 1 webcam-image-capture.png
   25  ffmpeg -hide_banner -f video4linux2 -list_formats all -i /dev/video0
   26  v4l2-ctl --list-devices
   27  ffmpeg -hide_banner -f video4linux2 -list_formats all -i /dev/video0
   28  shutdown
   29  exit
   30  apt update && apt install -y streamer
   31  streamer -f jpeg -o image.jpeg
   32  ffmpeg -f video4linux2 -list_formats all -i /dev/video0
   33  lsusb 
   34  v4l2-ctl --list-devices
   35  v4l2-ctl --list-devices --all

apt-get install fswebcam
fswebcam -r 640x480 --jpeg 85 -D 1 web-cam-shot.jpg

