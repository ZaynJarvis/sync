import React, { useRef, useState } from 'react';
import './file.css';

const FileInput = ({ setRefresh }) => {
  const fileInput = useRef(null);
  const [fileContents, setFileContents] = useState('');

  const handleFileSelect = (e) => {
    const file = e.target.files[0];
    const reader = new FileReader();

    reader.onload = (event) => {
      setFileContents(event.target.result);
      fetch('http://0.0.0.0:8888/sync', {
        method: 'POST',
        body: event.target.result,
      })
        .then(response => response.json())
        .then(() => setRefresh(Math.random()))
    };

    reader.readAsText(file);
  };

  return (
    <div style={{'marginBottom': '2rem'}}>
      <input type="file" id="videoFile" ref={fileInput} onChange={handleFileSelect} hidden/>
      <label for="videoFile">Upload Video</label>
      {fileContents ? <p className='preview'>{fileContents}</p> : <></>}
    </div>
  );
};

export default FileInput;