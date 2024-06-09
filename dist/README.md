# Spring Boot Application Monitor

## Requirements
- Operating System: Linux

## Description
This application is used to monitor health and metrics of Spring Boot applications with the use of actuator. Application is intended to run as a daemon service which will actively monitor with the provided monitor interval. 

There are 2 types of supported monitoring alerts.

1. Desktop Notifications/Alerts  
   Sends notifications about the CPU usage and JVM memory usage of applications. Sends alert notifications for critical health issues and when usage thresholds are exceeded. 

2. Email Alerts  
   Sends email alerts for critical health issues and when usage thresholds are exceeded.

## Usage

### Setup Configurations
- Add the desired configurations in the **.env** file inside **config** directory

### Install Monitoring Service
`sudo make install` 

### Check Service Status
`systemctl --user status springboot-app-monitor`

### Updating Configurations
- Change configurations of the **.env** inside **/usr/local/bin/springboot-app-monitor/config** directory
- Reload the daemon: `systemctl --user daemon-reload`

### Uninstalling Monitoring Service
`sudo make uninstall`