## Introduction

This is a simple Golang wrapper around the Splunk App Inspect REST-API, used for certifying the apps before submitting to the Splunkbase. 

This is an experimental version. Use it at your own risk. Tested on Mac OS X 10.14.6

A sample app file is provided. Please create an account at splunk.com if you dont have one. More info about the REST-API is available here : https://dev.splunk.com/enterprise/docs/releaseapps/appinspect/splunkappinspectapi/runappinspectrequestsapi

## Tool

The tool :
1) Generates an authentication token for the user
2) Uses the token and submits the Splunk App over REST-API for App Inspect Certification
3) Monitors the inspection status
4) Provides the raw results of the Inpection process

## Usage

```~/throw/goSplunk $ ./goSplunk -h
Usage of ./goSplunk:
  -f string
    	Splunk app file path
  -p string
    	Splunk password
  -u string
    	Splunk username
 ```
 
Example

```
./goSplunk -u vshesha -p mystery -f hashgen-1.0.0.tar.gz > file.txt
```
