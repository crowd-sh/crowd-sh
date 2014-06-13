# Stop the existing instance
sudo docker ps | grep abhiyerra/wmindex | awk '{print $1}' | xargs sudo docker kill

# Start the new instance
sudo docker run -d -p 127.0.0.1:3000:3000 abhiyerra/wmindex
