import Clappr from "@clappr/core";
import HlsjsPlayback from "@clappr/hlsjs-playback";
import { MediaControl, SpinnerThreeBounce } from "@clappr/plugins";

document.addEventListener('DOMContentLoaded', function() {
  if (!cfg || !cfg.broadcaster) {
    return;
  }

  const player = new Clappr.Player({
    source: cfg.staticHost+"/hls/"+cfg.broadcaster+".mp4/master.m3u8",
    plugins: [HlsjsPlayback, MediaControl, SpinnerThreeBounce],
    width: '100%',
    mute: true,
    parentId: "#player",
    poster: 'https://st.saverbate.com/images/performers/'+cfg.broadcaster+'_big.jpg',
    playback: {
      hlsjsConfig: {
        debug: false,
        maxBufferSize: 120*1000*1000,
        manifestLoadingTimeOut: 20000,
        manifestLoadingMaxRetry: 5,
        startFragPrefetch: true,
        manifestLoadingRetryDelay: 2000,
      }
    }
  });

  player.core.mediaControl.enable();
});
