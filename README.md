# WorkMachine

[![Build Status](https://travis-ci.org/workmachine/workmachine.png?branch=master)](https://travis-ci.org/abhiyerra/workmachine)

WorkMachine allows you to run crowdsourced workflows involving human
intervention. For example, take image transcription can be processed
using an OCR engine but that doesn't assure use that the quality will
be good so we can send it to a human to verify. The human input can
then be used to train the OCR algorithm. The engine is generic enough
to define arbitrary workflows which require human workers.

## Usage

### Create Workflow

curl -X POST -H "Content-Type: application/json" -d @british_library_tag.json localhost:3002/v1/workflows

### Create New Task

curl -X POST -H "Content-Type: application/json" -d @sample.json localhost:3002/v1/workflows/5/tasks

## Documentation

Check out the wiki for more information.

     https://github.com/workmachine/workmachine/wiki

## The MIT License (MIT)

Copyright (c) 2012-2015 Abhi Yerra <abhi@berkeley.edu>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
