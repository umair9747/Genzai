<h1 align="center">Genzai</h1>
<p align="center"><b>The IoT Security Toolkit</b></p>
<p align="center">
<a href="#description">Description</a> • <a href="#features">Features</a> • <a href="./docs/documentation.md#setup-usage">Setup & Usage</a> • <a href="#acknowledgements">Acknowledgements</a> • <a href="#contact">Contact</a><br>

<p align="center">
  <img src="https://img.shields.io/badge/Version-2.0-green">
  <img src="https://img.shields.io/badge/Black%20Hat%20Arsenal-%20Asia%202024-blue">
  <img src="https://img.shields.io/badge/GISEC Armory-%20Dubai%202024-blue">
  <a href="https://www.buymeacoffee.com/umair9747" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 21px !important;width: 94px !important;" ></a>
</p>

</p>
<hr>
<img src="./genzai.png">
<hr style="width:300px; height: 1px; margin: auto; margin-top: 20px;" />
<br>
<div id="description">
<h2> Description </h2>
Genzai is a powerful IoT security toolkit designed to identify and scan IoT devices for vulnerabilities and default credentials. It helps you identify IoT or Internet of Things related dashboards across a single or set of targets provided as input and furthermore scan them for default password issues and potential vulnerabilities based on paths and versions.
<br><br>
With its latest update, Genzai now offers an API-based approach for more flexible and scalable security assessments, supporting concurrent scanning and advanced target input formats.
<br><br>
Genzai currently supports fingerprinting over 20 IoT-based dashboards and has the same amount of templates to look for default password issues across them. It currently has a total of 10 vulnerability templates which will increase with coming updates.
</div>
<hr style="height: 1px;">

<div id="features">
<h2> Features </h2>

<h4>Fingerprinting - The Wappalyzer of IoT Devices</h4>
With Genzai, you can fingerprint the IoT Product running over a target based on the HTTP response received through it. Genzai can look for categories such as:

- Wireless Router
- Surveillance Camera
- HMI or Human Machine Interface
- Smart Power Control
- Building Access Control System
- Climate Control
- Industrial Automation
- Home Automation
- Water Treatment System

<h4>Default Password Checks</h4>
Based on the IoT product identified and the presence of a relevant template in the Vendor Logins DB, Genzai will check if the target is still using a vendor-specific default password.

<h4>Vulnerability Scanning</h4>
Also based on the IoT product identified and with the presence of a relevant template in the Vendor Vulns DB, Genzai will check for any potential vulnerabilities across the target.

<h4>New Features in Version 2.0</h4>

- RESTful API for easy integration with other tools
- Support for IP ranges and CIDR notation
- Concurrent scanning with customizable worker pool
- Verbose logging option for detailed scan information
- Improved error handling and reporting
</div>

<div id="acknowledgements">
<h2> Acknowledgements </h2>
Genzai has been or will be noticed at:
<ul type="disc">
<li><a href="https://www.blackhat.com/asia-24/arsenal/schedule/index.html#genzai---the-iot-security-toolkit-37373">Black Hat Asia 2024 [Arsenal]</a></li>
<li><a href="https://www.gisec.ae/gisec-armory">GISEC Armory Edition 1 Dubai 2024</a></li>
</ul>

Special thanks to:
<ul type="disc">
<li><a href="https://github.com/rumble773">rumble773</a> for significant contributions to the API implementation and concurrent scanning features.</li>
</ul>
</div>

<div id="contact">
<h2> Contact </h2>
If you have any questions or feedback about Genzai, feel free to reach out:

- Umair Nehri (Original Author): <a href="https://in.linkedin.com/in/umair-nehri-49699317a">LinkedIn</a> or <a href="mailto:umairnehri9747@gmail.com">Email</a>
- rumble773 (Major Contributor): <a href="https://github.com/rumble773">GitHub</a>
</div>

<h2>Legal Disclaimer</h2>
Usage of Genzai for scanning or attacking targets without prior mutual consent is illegal. It is the end user's responsibility to obey all applicable local, state and federal laws. Developers assume no liability and are not responsible for any misuse or damage caused by this program.