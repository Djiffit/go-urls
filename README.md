## Great link store application

This basic link store application can create links with aliases! When these aliases are fetched an event is logged into visit table. Thus this can be extremely useful piece of software! The application also has users that you can create and log in with. Of course the important funcitonality requires the use of the jwt token in the authorization header. To get a review of all links in the application, you can go to the main links page and find the top links, their total visits, top visitor and how many times the top visitor has visited the given link.

So basically just testing creating basic Go REST api with jwt authorization and psql tables.

Of course when nobody wants to use this you will need a Postgres database running somewhere and configure the constants in the main file, I would suggest running it in Docker but do what you will, imaginary friend!
