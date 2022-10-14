<h1 align="center">Kicking Off Hacktoberfest with ACM-VIT!</h1>
<p align="center">
<img src="./banner.png">
</p>

<h2 align="center"> Golly Bot! </h2>

<p align="center"> 
A general purpose Discord bot, written in GO!
</p>

<p align="center">
  <a href="https://acmvit.in/" target="_blank">
    <img alt="made-by-acm" src="https://img.shields.io/badge/MADE%20BY-ACM%20VIT-blue?style=for-the-badge" />
  </a>
    <!-- Uncomment the below line to add the license badge. Make sure the right license badge is reflected. -->
    <!-- <img alt="license" src="https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge" /> -->
    <!-- forks/stars/tech stack in the form of badges from https://shields.io/ -->
</p>

---
## Submitting a Pull Request

 * Fork the repository by clicking the fork button on top right corner of the page
 * Clone the target repository. To clone, click on the clone button and copy the https address. Then run 
 <pre><code>git clone https://github.com/ACM-VIT/golly-bot.git</code></pre>
* Go to the cloned directory by running 
<pre><code>cd golly-bot</code></pre>
* Create a new branch. Use 
<pre><code> git checkout -b mynewbranch</code></pre>
* Make your changes to the code. Add changes to your branch by using 
<pre><code>git add .</code></pre>
* Commit the chanes by executing
<pre><code>git commit -m "short message describing changes"</code></pre>
* Push to remote. To do this, run 
<pre><code>git push origin mynewbranch</code></pre>
* Create a pull request. Go to the target repository and click on the "Compare & pull request" button. **Make sure your PR description mentions which issues you're solving.** Wait for your request to be accepted, and you're good to go!
<p align="center"><img src="https://drive.google.com/u/1/uc?id=1f9JKAR-kRvCRGxIs_SAvegaYDPx53T9G&export=download"></img></p>

---
## Guidelines for Pull Request

<!-- general guidelines here -->
  * Avoid pull requests that :
      * are automated or scripted
      * that are plagarized from someone else's branch
  * Do not spam
  * Project maintainer's decision on validity of PR is final.

  For additional guidelines, refer to [this website](https://hacktoberfest.com/participation/) for additional information

---

## Instructions on how to run the bot locally
 
  * In the cloned directory :
     * Add a .env file following the .env.sample file provided and add the required secrets like token,etc. <br>(Instructions on how to create a Discord Bot token can be found [here](https://www.writebots.com/discord-bot-token/))
     * Allow the bot to have Privilaged Intents as found [here](https://support.discord.com/hc/en-us/articles/360040720412-Bot-Verification-and-Data-Whitelisting)
     * Make sure you have go and all its dependencies installed in your system. 
       <br>(Instructions can be found [here](https://go.dev/doc/install))
     * Run the following command
       <pre><code>go run main.go</code></pre>
## Instructions on how to run the bot in a docker container

  * Build the docker image from the Dockerfile using the following command, **make sure you are inside the project folder.**
    <pre><code>docker build -t &lt;your image name&gt; .</pre></code>
  * Now run the application inside the container using the following command
    <pre><code>docker run -it &lt;your image name&gt;</pre></code>
  * To specify command line arguments, use
    <pre><code>docker run -it &lt;your image name&gt; &lt;arguments&gt;</pre></code> 
    For example:<pre><code>docker run -it &lt;your image name&gt; -rmcmd=false</pre></code>
  
<!-- ---
## Overview

The overview starts here. Random text about the project, motive, how, what, why etc.

---
## Usage
How To, Features, Installation etc. as subheadings in this section. example

Lets get started!
```console
git remote add
git fetch
git merge
``` -->

---
## License
[License](LICENSE)

---
## Authors

  - [Anish Raghavendra](https://github.com/z404)
  - [Manav Muthanna](https://github.com/ManavMuthanna)
  - [Sarthak Gupta](https://github.com/gptsarthak)  
<!-- **Contributors:** Generate contributors list using this link - https://contributors-img.web.app/preview -->
