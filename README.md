# Spring Boot Application Monitor

## Description
This application is used to monitor health and metrics of Spring Boot applications with the use of actuator. Application is intended to run as a daemon service which will actively monitor with the provided monitor interval. 

There are 2 types of supported monitoring alerts.

1. Desktop Notifications/Alerts  
   Sends notifications about the CPU usage and JVM memory usage of applications. Sends alert notifications for critical health issues and when usage thresholds are exceeded. 

2. Email Alerts  
   Sends email alerts for critical health issues and when usage thresholds are exceeded.