Contributing
------------
Golossus is an open source, community-driven project.

If you'd like to contribute, please read the following sections.

Submitting a patch
------------------
A pull request is the best way to propose changes to Golossus. Follow next steps to propose yours.

### Step 1: Check existing Issues and Pull Requests
Before start any work, check to see if someone else also stated the same problem or even started working on a pull 
request by searching on [GitHub][1].

If you are unsure or have any questions, please ask your questions on the #contribs channel on [Golossus Slack][2].

### Step 2: Setup your Environment
You will need to have **Git** installed and an account on Github. Set up your user information to match your Github account 
details, but use your real name:

```bash
git config [--global or --local] user.name "Your Name"
git config [--global or --local] user.email your.github@email.com
```

Fork the [Golossus repository][1] to your Github account and ensure current tests pass before working on your Pull Request.

### Step 3: Work on your Pull Request
Before you start, you should be aware that all the code you are going to submit must be released under the [MIT license][3].

Work on current stable branch for fixes, and on master for new features. [Golossus team][4] will decide if some changes must 
be propagated to any other branches. 

You must create a Topic Branch from the previous ones, and provide a descriptive name including the ticket number if an
issue for the code you're about to write exists.

Work on the code as much as you want and commit as much as you want, but keep in mind to follow the rules on this document.

Add unit tests to prove that the bug is fixed or that the new feature actually works.

Try hard to not break backward compatibility, if you must do so, take into account that Pull Requests that break backward 
compatibility have less chance to be merged.

Do atomic and logically separate commits (use the power of git rebase to have a clean and logical history).

Never fix coding standards in some existing code as it makes the code review more difficult.

Write good commit messages.

### Step 4:Prepare your Pull Request for Submission
When your Pull Request is not about a bug fix (when you add a new feature or change an existing one for instance), it 
must also include the following:

* An explanation of the changes in the relevant CHANGELOG file(s), use the [BC BREAK] or the [DEPRECATION] prefix when relevant).
* An explanation on how to upgrade an existing application in the relevant UPGRADE file(s) if the changes break backward 
  compatibility or if you deprecate something that will ultimately break backward compatibility.
* Add your name into the CONTRIBUTORS file if you're not yet added.

### Step 5: Submit your Pull Request
Whenever you feel that your Pull Request is ready for submission, follow the following steps.

* Rebase your Pull Request, Before submitting it.
* Fix merge conflicts if they appear during rebase.
* Check that all tests still pass and push your branch remotely.
* Finally make a pull request on the Golossus GitHub repository.

### Step 6: Receiving Feedback
We ask all contributors to follow a constructive feedback process.

If you think someone fails to keep this advice in mind and you want another perspective, please join the #contribs 
channel on [Golossus Slack][2]. If you receive feedback you find abusive please contact the [conduct][5].

The Golossus team is responsible for deciding which code gets merged, so their feedback is the most relevant. 
So do not feel pressured to refactor your code immediately when someone provides feedback.

Based on the feedback on the pull request, you might need to rework your Pull Request. As in step 5, you will need to
rebase again and ensure tests are passing.

Golossus team will ask you to "squash" your commits before merging it. This means you will convert many commits to one commit. 

Reporting a bug
---------------

If you find a bug in a Golossus module, please report it and help us make a better Golossus.

> CAUTION
>
> If you think you've found a security issue, please use the special [procedure][6] instead.

Before submitting a bug:

* Double-check the official documentation to see if it is just a misusing topic.
* Ask for assistance on Stack Overflow and/or on the #support channel of the [Golossus Slack][2].
* If your problem definitely looks like a bug, report it using the official bug tracker of the corresponding package on   
  Github. 

Basic rules to follow:

* Use the title field to clearly describe the issue.
* Describe the steps needed to reproduce the bug with short code examples. It's better to provide a test.
* Give as much detail as possible about your environment (OS, Go version, Golossus modules version, ...).
* Any log, output or trace you provide must be done in plain text. Do not provide it as a screenshot, and be careful with
  sensitive information you might be sharing by mistake. 
* Optionally attach a patch.

Reporting a security issue
--------------------------

If you think that you have found a security issue in a Golossus package, don't use the bug tracker and don't make it 
public. Instead, all security issues must be sent to [security@golossus.com][7]. Emails sent to this address are 
forwarded to the Golossus team.

### Resolution process
> While the Golossus team is working on a patch, please do not reveal the issue publicly. We will resolve the problem as
> fast as possible depending on its severity.

For each report, we first try to confirm the vulnerability. When it is confirmed, the Golossus team  will work on a solution
and they will follow the next series of steps:

1. Notify the reporter to acknowledge the vulnerability.
2. Work hard on a patch.
3. Get a CVE identifier from [mitre.org][8].
4. Write a security post announcing the vulnerability for the [Golossus blog][9]. 
5. Send the patch and the contents of the announcement to the reporter for review.
6. Apply the patch and package a new version.
7. Publish the post and add an entrance to the [security advisories list][10]

### Dependant Open-Source projects collaboration
If you are using Golossus for your own open-source or project, let us know by sending an email to [downstream@golossus.com][13]. 
Include information about your project and/or company, and which Golossus component your are using. After verification of 
the information provided your company and/or project will be added into our security notification mailing list. Do not 
forget to include at least one email address where you want to receive security notifications.

After notifying a security issue to the dependant projects, we will try to approach the right people of the project to set
a strategy to help to update and adopt the fixes. 

> By sending the security issue notification request email, and after its approval, you allow us to contact by means of 
> the email(s) address(es) provided, not only to get notified about a security issue but also to ask for permission to 
> publish your project/company as a downstream project using Golossus.

#### Projects using Golossus

Currently there's no know project using actively Golossus, we expect this to change in the future.

### Severity scoring

In order to determine the severity of a security issue we take into account:
* The complexity of any potential attack. 
* The impact of the vulnerability.
* The number of projects that may be affected. 

The Golossus team has not set yet a quantitative mechanism to score the severity of an issue, but taking into account the
points mentioned before, we will categorize an issue into a level of: Low, Medium, High, or Critical.

### Security advisories

Check the Security Advisories blog category for a [list of all security vulnerabilities][10] 
that were fixed in previous releases.

Golossus team
-------------

The Golossus team is the group of developers that drives and decides the direction and evolution of the Golossus project. 
They have the power to reject, accept and merge the features and patches proposed by the community.

This section states the rules Golossus members must adhere and follow. These rules are effective upon publication of this 
document. 

### Leaders
* Santiago Garcia (sangarbe)
* Jorge Brisa (jorbriib)

### Former members
Thanks for your contributions even if you're no longer member of the team:
* David Camarena (davidcamarena)

### Membership application
All the Golossus team members should be long-time contributors with demonstrated knowledge and commitment to push the 
project forward. Member applications must be sent to [members@golossus.com][11]. After evaluation 
by the team leaders, needed privileges will be granted to any accepted applications' sender.

### Membership revocation
Any member of the Golossus team may loose their membership if one of the following circumstances apply:
* Not following the rules, conventions or code of conduct.
* Long-time inactivity.
* Demonstrated negligence or intention to harm the project.
* Leader(s) decision.

Development rules
-----------------

Golossus project development is based on pull requests proposed by any member of the community. However, this project 
is young and pull request acceptance or rejection is decided based on the criteria of Golossus team members.

### Coding standards and conventions

Golossus does not supply its own set of coding standards, we stick to conventions directly addressed by the [go official
documentation][12]. Let us know if we don't follow them properly.

### Merging process
All code must be committed to the repository through pull requests. After enough time for reviews by the team members, a
pull request can be merged in the following cases:
* Community pull request: Every team leader must accept the code.
* Team member pull request: Only one team leader must accept the code.
* Team leader pull request: A different team leader must accept the code.

Only the team leaders can eventually merge a pull request. 

Even if team members (except team leaders) have not enough privileges to merge a pull request, they are allowed to ask 
for code completeness, refinement or changes in community pull request. Their help and effort is highly valuable to 
accelerate the lead time to accept and/or reject a pull request. 

Regardless team members approval/rejection, team leaders may change the final decision.

### Release policy
Team leaders are in charge of releasing every Golossus package version.

Amendments
----------
The rules described in this document may change at anytime at the discretion of the project leaders.
 
[1]: https://github.com/golossus  
[2]: https://join.slack.com/t/golossus/shared_invite/zt-db4brnes-M8q1Lw2ouFT5X~gQg69NQQ
[3]: ./LICENSE
[4]: #golossus-team
[5]: mailto:conduct@golossus.com
[6]: #reporting-a-security-issue
[7]: mailto:security@golossus.com
[8]: https://cve.mitre.org/cve/request_id.html
[9]: http://www.golossus.com/blog
[10]: http://www.golossus.com/security-advisories
[11]: mailto:members@golossus.com
[12]: https://golang.org/doc/
[13]: mailto:downstream@golossus.com
