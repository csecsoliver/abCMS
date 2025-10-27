# abCMS

Official instance: https://abcms.olio.ovh

## What's new? (abCMS.cozy)

You can now post images on abCMS.cozy! It's a fall themed photo sharing sister-site to abCMS.text with the same base philosophy. The mana system has not been implemented on the cozy site and might never be. Depends on my mood basically.

There are two main ways you can upload photos:
1. Through the website: You can just select the image, give it a title (optional) and post it. (**Beware, no delete function implemented yet!**)
2. Through the [Custom Uploader](Shttps://github.com/SrS2225a/custom_uploader) ([F-droid](https://f-droid.org/en/packages/com.nyx.custom_uploader/), [Google Play](https://play.google.com/store/apps/details?id=com.nyx.custom_uploader)):
    - Log into the abCMS.cozy site in a browser an set an upload key, this will be you key to upload files without publishing them yet.
    - In the app, go to settings and add an uploader
    - Set the url to ```https://abcms.olio.ovh/cozy/postimg?token=yourkey&user=youruser```
    - Of course replace the url with your server of choice, ```yourkey``` with your chosen key and ```youruser``` with your username on the site
    - And set the Form Data Name to ```image```
    - File encoding can be left off
    - Save this uploader and then tap on it in the list to select it. From now on you can upload images in the app, or share them from your gallery with the app.
    - If you have uploaded some photos this way, you can log in to the website and go through your uploads, filtering them out if needed.

The endpoint for the app uploader can accept multipart/form-data with one field named ```image``` from anywhere, so feel free to use your own solution.

The official instance has plenty storage, so feel free to post your best images taken this time of year.

## abCMS.text

The abc management system. It will be an all in one system specifcally made to manage my homelab and related services.  
Right now however it is a bit off the tracks, it has become a blog platform encouraging frequent, but short posts, rewarding them with mana!!  

- You can sign up in the Dashboard with the instance's signup secret, **which is basically a key to the instance, so hosters can close it down a bit.**   
- It is usually "securesignup" as specified in the code, but instance hosters can specify their own in a .env file.  
> However, the official instance hosted by me has it set to an empty string, in which case it is not required to enter it at signup. 
- After signing up you can make posts, but posts cannot be edited or deleted to emphasize the quick, short posts nature of the site.  
- You earn mana for every post, and you can spend it on enchanting your posts, or just brag about how much you have.

Another use I found:
### Do you miss the devlogs from SOM?
Well abcms is here to fill that void! Simply make an account on https://abcms.olio.ovh (the url it redirects to changes sometimes, so bookmark this one) and start posting and scroling away.

******

# My experience  
This is my first time working with htmx, but not my first time in the bottle.js ecosystem. I am not proficient in either of them right now, but I am getting used to their quirks.  
The htmx seems to be most like the php philosophy, which is not exactly what bottle is made for. I have only used very basic templates, using the python builtin .format() method, but the template() function in bottle seems intriguing.  
This is also my very first time using jwt authentication, and I already like it a lot more than I did my previous attepmts at token management.  
I also came to hate browser caching altogether, which is unreasonable, but it cost me like an hour of on-off debugging to figure out why the server did not log my request, despite the content appearing on the site.  

(Yes, I did rebrand the app to cover up my lack of time to implement stuff)
