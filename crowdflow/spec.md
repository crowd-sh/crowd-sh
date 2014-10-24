Assignments should be done regardless of the backend. It doesn't
matter what the backend is. It can be the website, mobile, mturk but
it just needs to be able to find the correct task at the right time
without the library writer having t odeal with it.

This should mean that there is a generic interface with a Serve
interface passed.

This would create the appropriate assignments for the correct view. So
if a user requests on mobile it'll split an assignment and give it to
the mobile user.

If it is a MTurk user it passes the whole assignment to the user. If
it is a regular user it passes the whole page.
