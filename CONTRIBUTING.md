Branching and Releasing
=========

**TL;DR** Create your Pull Requests targeting `dev` branch!

We are using git tags for releases. Version numbers doesn't have any suffix/prefix, for example: **6.2.12**

New features should be tested by developers in `sysnet` branch which represents an isolated subutai ecosystem and can be used by project maintainers for integration tests.

After isolated tests code from `sysnet` branch must be merged into `dev` branch. In a dev branch developers, contributors and project maintainers can test p2p using real-life applications and environments.  

`master` branch should be stable, containing all the hotfixes along with the new features that was tested and merged from `dev` branch. Master is used by Subutai QA team for their tests. 

Once in a while code in `master` is being tagged.

Coding style and tests
===========

Try to give clean names to all your variables, functions and methods, structs, types and everything else. Do not forget to use linter. 

And always try to write a test function for your code!

Contributors
============

Contributors will be listed in alphabet order. If we forgot to mention you - please contact us directly. 

* [chennqqi](https://github.com/chennqqi)