<h1 align="center">Genzai</h1>
<p align="center"><b>The IoT Security Toolkit</b></p>
<p align="center">
<a href="../README.md#description">Description</a> • <a href="../README.md#features">Features</a> • <a href="#setupnusage">Setup & Usage</a> • <a href="../README.md#acknowledgements">Acknowledgements</a> • <a href="../README.md#contact">Contact Me</a><br>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/Version-1.0-green">
  <img src="https://img.shields.io/badge/Black%20Hat%20Arsenal-%20Asia%202024-blue">
  <img src="https://img.shields.io/badge/GISEC Armory-%20Dubai%202024-blue">
  <a href="https://www.buymeacoffee.com/umair9747" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 21px !important;width: 94px !important;" ></a>
</p>

<hr>


<div id="setupnusage">
<h2> Setup and Usage </h2>
<h3> Setup </h3>
The tool is written in Go, so make sure to <a href="https://go.dev/dl/">install</a> it on your system before proceeding. The setup is pretty easy and straight forward. Just follow the below steps in order to quickly install and get the binary working.
<br>
<br>
Firstly clone the repo to any directory/path of your liking,<br><br>

```
git clone https://github.com/umair9747/Genzai.git
```
Afer this, just run the following command in order to build the binary according to your environment.

```
go build
```

<h3> Usage </h3>

<h4>1. Basic Usage</h4>
In order to get started with Genzai and run it straightaway, you just need to provide your target(s) as input to the tool. This can be mainly done in the following ways,

<br><br>
To run it against a single target,

```
./genzai http://1.1.1.1/
```
To run it against multiple targets passed directly through CLI,

```
./genzai http://1.1.1.1/ http://8.8.8.8/
```
To run it against multiple targets passed through an input file,

```
./genzai targets.txt
```

<h4>2. Output</h4>
If you'd like to log the output from Genzai in some file, you have the <b>-save</b> flag to the rescue! The flag will log the output in a file named output.json by default if no additional value is provided along with the flag.

<br><br>
So, in order to log the output in a specific file name, you will use,

```
./genzai http://1.1.1.1 -save myfile.json
```

And with the below example, it will be saved to output.json by default.
```
./genzai http://1.1.1.1 -save
```

<br><br>
Let's also discuss about the format of output that Genzai returns. It will be in the below format:

<br>

```
{
  "Response": {
    "Results": [
      {
        "Target": "",
        "IoTidentified": "",
        "category": "",
        "Issues": [
          {
            "IssueTitle": "",
            "URL": "",
            "AdditionalContext": ""
          }
        ]
      }
    ],
    "Targets": []
  }
}

```
The results array contains the entries for all the IoT related assets that were successfully identified and scanned. Target field contains the URL of the asset, IoTidentified contains the product name, category contains the exact category the IoT product belongs to, Issues array will be populated with all the issues identified with the asset such as any potential vulnerabilities and default password issues.<br>
Finally, the Targets array contains the list of all targets that were scanned using the tool irrespective of them being identified as an IoT asset or not.


The below example output would hopefully give you a glimpse of the format,

<br>

```
./genzai http://1.1.1.1/

::::::::   :::::::::: ::::    ::: :::::::::     :::     ::::::::::: 
:+:    :+: :+:        :+:+:   :+:      :+:    :+: :+:       :+:     
+:+        +:+        :+:+:+  +:+     +:+    +:+   +:+      +:+     
:#:        +#++:++#   +#+ +:+ +#+    +#+    +#++:++#++:     +#+     
+#+   +#+# +#+        +#+  +#+#+#   +#+     +#+     +#+     +#+     
#+#    #+# #+#        #+#   #+#+#  #+#      #+#     #+#     #+#     
 ########  ########## ###    #### ######### ###     ### ########### 

        The IoT Security Toolkit by Umair Nehri (0x9747)


2024/03/30 23:19:47 Genzai is starting...
2024/03/30 23:19:47 Loading Genzai Signatures DB...
2024/03/30 23:19:47 Loading Vendor Passwords DB...
2024/03/30 23:19:47 Loading Vendor Vulnerabilities DB...

 

2024/03/30 23:19:47 Starting the scan for http://1.1.1.1/
2024/03/30 23:19:49 IoT Dashboard Discovered: TP-Link Wireless Router
2024/03/30 23:19:49 Trying for default vendor-specific [ TP-Link Wireless Router ] passwords...
2024/03/30 23:19:51 http://1.1.1.1/ [ TP-Link Wireless Router ] is vulnerable with default password -  TP-Link Router Default Password - admin:admin
2024/03/30 23:19:51 Scanning for any known vulnerabilities from the DB related to TP-Link Wireless Router
2024/03/30 23:19:57 http://1.1.1.1/ [ TP-Link Wireless Router ] is vulnerable  -  TP-LINK Wireless N Router WR841N Potentially Vulnerable to Buffer Overflow - CVE-2020-8423

2024/03/30 23:20:45 No file name detected to log the output. Skipping to printing it!

 
{
    "Results": [
        {
            "Target": "http://1.1.1.1/",
            "IoTidentified": "TP-Link Wireless Router",
            "category": "Router",
            "Issues": [
                {
                    "IssueTitle": "TP-Link Router Default Password - admin:admin",
                    "URL": "http://1.1.1.1/userRpm/LoginRpm.htm?Save=Save",
                    "AdditionalContext": "The resulting body had matching strings from the DB."
                },
                {
                    "IssueTitle": "TP-LINK Wireless N Router WR841N Potentially Vulnerable to Buffer Overflow - CVE-2020-8423",
                    "URL": "http://1.1.1.1/",
                    "AdditionalContext": "The resulting headers matched with those in the DB."
                }
            ]
        }
    ],
    "Targets": [
        "http://1.1.1.1/"
    ]
}
```
</div>