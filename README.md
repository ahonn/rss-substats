# rss-substats
Count How Many People are Subscribed to Your RSS, based on [Substats](https://github.com/spencerwooo/Substats)

You can count the number of people subscribed to you via `https://api.spencerwoo.com/substats?source={SOURCE}&queryKey={QUERY_KEY}`, but Substats has poor support for counting multiple Query Keys. This project combines the number of RSS subscribers by combining multiple substats requests.

# Useage

`https://rss-substats.ahonn.me/?source=feedly,inoreader&queryKey=https://www.ahonn.me/atom.xml,https://www.ahonn.me/rss.xml`

Combine the number of subscriptions of the two query keys corresponding to feedly and inoreader

And you'll see that:

```json
{
  "status":200,
  "data":{
    "totalSubs":209,
    "subsInEachSource":{
      "feedly":111,
      "inoreader":98
    },
  "failedSources":{}
  }
}
```

# License
MIT License
