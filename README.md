# abCMS
The abc management system. It will be an all in one system specifcally made to manage my homelab and related services.  
Right now however it is a bit off the tracks, it has become a blog platform encouraging frequent, but short posts, rewarding them with coins!!  

- You can sign up in the Dashboard with a signup secret, which is usually "securesignup" as specified in the code, but instance hosters can specify their own in a .env file.  
- This secret is basically a key to the instance, so hosters can close it down a bit.  
- After signing up you can make posts, but posts cannot be edited or deleted to emphasize the quick, short posts nature of the site.  
- You earn coins for every post, and you can brag with those if you want. B)


# My experience  
This is my first time working with htmx, but not my first time in the bottle.js ecosystem. I am not proficient in either of them right now, but I am getting used to their quirks.  
The htmx seems to be most like the php philosophy, which is not exactly what bottle is made for. I have only used very basic templates, using the python builtin .format() method, but the template() function in bottle seems intriguing.  
This is also my very first time using jwt authentication, and I already like it a lot more than I did my previous attepmts at token management.  
I also came to hate browser caching altogether, which is unreasonable, but it cost me like an hour of on-off debugging to figure out why the server did not log my request, despite the content appearing on the site.  

(Yes, I did rebrand the app to cover up my lack of time to implement stuff)
