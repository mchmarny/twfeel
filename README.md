# twfeel

On-demand Twitter topic sentiment provider

## Demo (WIP)

* Deploy service to Knative cluster
* Configure chat bot

### Knative


### Chat Bot

In the Google API Console, enable the Hangouts Chat API by doing the following:

* In the navigation, click APIs & Services > Dashboard.
* In the Dashboard, click Enable APIs and Services.
* Search for "Hangouts Chat API" and enable the API.

Once the API is enabled, click the Configuration tab. In the Configuration pane, do the following:

* In the Bot name field, enter 'YOUR_NAME'.
* In the Avatar URL field, enter 'https://goo.gl/yKKjbw'.
* In the Description field, enter 'Tweet sentimenter'.
* Under Functionality, select Bot works in direct messages and room
* Under Connection settings, select Bot URL and paste the URL for the Knative service
* Set the Permissions for all in the room

When you've finished configuring your bot, click Save Changes.


