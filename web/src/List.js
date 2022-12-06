import React, { useState, useEffect } from 'react';
// import Player from 'xgplayer';
import VideoJS from './video'
import 'video.js/dist/video-js.css';
import './App.css';

const List = () => {
  // Initialize the state with an empty array
  const [data, setData] = useState([]);
  const [video, setVideo] = useState("");


  // Use the useEffect hook to make the API call and update the state
  useEffect(() => {
    fetch('http://0.0.0.0:8888/vid')
      .then(response => response.json())
      .then(data => setData(data))
  }, []);

  // Initialize the state with an empty array
  const [play, setPlay] = useState("");


  // Use the useEffect hook to make the API call and update the state
  useEffect(() => {
    if (video) {
      fetch('http://0.0.0.0:8888/play/' + video)
        .then(response => response.json())
        .then(data => setPlay(data))
    }
  }, [video]);

  return (
    <div className='page'>
      <div className='list'>
        <ul style={{ 'margin': 0 }}>
          {data.map((i) =>
            <li key={i.ID} className={'status' + i.Status}>
              <code onClick={() => setVideo(i.ID)}
              // {'color': i.Status == 2 ? 'success' : i.Status == 1 ? 'pending' : 'failed'}
              >{i.ID}</code>
            </li>
          )}
        </ul>
      </div>
      <div>
        {video !== "" ? (<div >
          <p>{video}</p>
          {
            play.url &&
            <>
              <p>{play.reason}</p>
              <a target='_blank' href={play.url}>{play.type}</a>
              <VideoJS options={{
                autoplay: true,
                controls: true,
                responsive: true,
                fluid: true,
                sources: [{
                  src: play.url,
                  type: 'video/mp4'
                }]
              }} />
              {/* <p style={{ 'fontSize': '5px' }}>{play.url}</p> */}
            </>
          }
        </div>) : <></>}

      </div>
    </div>
  );
}

export default List;