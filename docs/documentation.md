<h1 align="center">Genzai</h1>
<p align="center"><b>The IoT Security Toolkit</b></p>
<p align="center">
<a href="../README.md#description">Description</a> • <a href="../README.md#features">Features</a> • <a href="#setupnusage">Setup & Usage</a> • <a href="#acknowledgements">Acknowledgements</a><br>
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

<br>
To run it against a single target,

```
./genzai http://1.1.1.1/
```
<br>
To run it against multiple targets passed directly through CLI,

```
./genzai http://1.1.1.1/ http://8.8.8.8/
```
<br>
To run it against multiple targets passed through an input file,

```
./genzai targets.txt
```
</div>