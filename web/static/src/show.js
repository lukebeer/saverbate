import Clappr from "@clappr/core";
import HlsjsPlayback from "@clappr/hlsjs-playback";
import { MediaControl, SpinnerThreeBounce } from "@clappr/plugins";

document.addEventListener('DOMContentLoaded', function() {
  if (!cfg || !cfg.broadcaster) {
    return;
  }

  const player = new Clappr.Player({
    source: "https://st.saverbate.com/hls/"+cfg.broadcaster+".mp4/master.m3u8",
    plugins: [HlsjsPlayback, MediaControl, SpinnerThreeBounce],
    mute: true,
    parentId: "#player",
    poster: 'http://imgproxy.saverbate.localhost/rs/fit/800/0/sm/0/plain/local://static/'+cfg.broadcaster+'.jpg',
    playback: {
      hlsjsConfig: {
        debug: false,
      }
    }
  });

  player.core.mediaControl.enable();
});
