Note: This is under active development, there are no tests, and behavior can be eratic.
We are working on making it better but until then here be demons.

# WorkMachine

WorkMachine allows you to run crowdsourced workflows involving human
intervention. For example, take image transcription can be processed
using an OCR engine but that doesn't assure use that the quality will
be good so we can send it to a human to verify. The human input can
then be used to train the OCR algorithm. The engine is generic enough
to define arbitrary workflows which require human workers.

## Usage

### Config

```
{
    "Title": "Find Information for Company",
    "Description": "Find the Company URL information for the information provided. Incorrect information will be REJECTED.",
    "Tags": "research,data,company,leads",
    "Reward": "0.03",
    "AirTable": {
        "Base": "basename",
        "API": "API",
        "Table": "Table"
    },
    "Fields": [
        {
            "Name": "Company Name",
            "Type": "ShortText",
            "Description": "Use the company name to fill out the fields below.",
        },
        {
            "Name": "Company Website",
            "Type": "ShortText",
            "Description": "Insert the company website if it doesn't exist. Find the company name using Google.",
            "Value": "",
            "Options": null
        },
        {
            "Name": "Company LinkedIn",
            "Type": "ShortText",
            "Description": "Use LinkedIn. Find the Company's LinkedIn url here. Put 'not found' if you cannot find the information. ",
            "Value": "",
            "Options": null
        }
    ]
}
```

### Staging
```
workmachine --review config.json
```

### Production

```
workmachine --live --review --output config.json
```

# Project by opsZero

<a href="https://www.opszero.com"><img src="http://assets.opszero.com.s3.amazonaws.com/images/opszero_11_29_2016.png" width="300px"/></a>

This project is brought to you by [opsZero](https://www.opszero.com) we
provide DevOps and Cloud Infrastructure as a Service for Startups. If you
need help with your infrastructure reach out.

## The MIT License (MIT)

Copyright (c) 2012-2017 opsZero <abhi@opszero.com>

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