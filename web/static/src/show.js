import Clappr from "@clappr/core";
import HlsjsPlayback from "@clappr/hlsjs-playback";
import { MediaControl, SpinnerThreeBounce } from "@clappr/plugins";

document.addEventListener('DOMContentLoaded', function() {
  const player = new Clappr.Player({
    source: "https://bitdash-a.akamaihd.net/content/sintel/hls/playlist.m3u8",
    plugins: [HlsjsPlayback, MediaControl, SpinnerThreeBounce],
    mute: true,
    parentId: "#player",
    poster: 'http://imgproxy.saverbate.localhost/rs/fit/800/0/sm/0/plain/local://static/jykfqy/a522691f-0402-4718-a180-a165db6bf0d8.jpg',
    playback: {
      hlsjsConfig: {
        debug: false,
      }
    }
  });

  player.core.mediaControl.enable();
});
