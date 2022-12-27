import React, { useState, useEffect } from 'react';
import VideoJS from './Video'
import FileInput from './FileInput'
import 'video.js/dist/video-js.css';
import './List.css';
import Link from './link.svg'

const List = () => {
  // Initialize the state with an empty array
  const [data, setData] = useState([]);
  const [video, setVideo] = useState("");
  const [vid, setVID] = useState("");
  const [refresh, setRefresh] = useState(0);


  // Use the useEffect hook to make the API call and update the state
  useEffect(() => {
    fetch('http://0.0.0.0:8888/vid')
      .then(response => response.json())
      .then(data => setData(data))
  }, [refresh]);

  // Initialize the state with an empty array
  const [play, setPlay] = useState("");


  // Use the useEffect hook to make the API call and update the state
  useEffect(() => {
    if (video) {
      fetch('http://0.0.0.0:8888/play/' + video)
        .then(response => response.json())
        .then(data => { setPlay(data); setVID(data.vid); setRefresh(Math.random()) })
    }
  }, [video]);

  return (
    <div className='page'>
      <div className='list'>

        <div className='id-container'>
          <FileInput setRefresh={setRefresh} />
          <ul className='horizontal'>
            <li className='status1'>Pending</li>
            <li className='status2'>Success</li>
            <li className='status3'>Failed</li>
          </ul>
          <ul className='veritical'>
            {data ? data.map((i) =>
              <li key={i.ID} className={'status' + i.Status}>
                <code className={video === i.ID ? 'selected' : ''} onClick={() => setVideo(i.ID)}
                >{i.ID}</code>
              </li>
            ) : <></>}
          </ul>
        </div>
      </div>
      <div className='video-container'>
        {video !== "" ? (<div >
          <p className='title'>{'Video ID: ' + video}</p>
          <p className='title'>{'VID: ' + vid}</p>
          {
            play.url &&
            <>
              <div className='badge'>
                <span className={play.type == 'source_url' ? 'url badge-pending' : 'url badge-good'}>
                  <a className='tooltip' target='_blank' href={play.url}>{play.type}
                    <span className='tooltiptext'>{play.url.split("?")[0]}</span>
                  </a>
                  <img src={Link} className="link" alt="link" />
                </span>
                {play.reason ? <span className='badge-warning'>{'reason: ' + play.reason}</span> : <></>}
              </div>
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
            </>
          }
        </div>) : <></>}

      </div>
    </div>
  );
}

export default List;