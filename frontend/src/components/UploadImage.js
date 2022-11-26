import React, { useRef, useState } from 'react';
import useAxiosPrivate from '../api/useAxiosPrivate';
import {
  FileUploadContainer,
  FormField,
  DragDropText,
  UploadFileBtn,
  FilePreviewContainer,
  ImagePreview,
  PreviewContainer,
  PreviewList,
  FileMetaData,
  RemoveFileIcon
} from '../styles/file-upload';
import '../styles/file-upload.css';

const KILO_BYTES_PER_BYTE = 1000;

/* const convertNestedObjectToArray = (nestedObj) =>
  Object.keys(nestedObj).map((key) => nestedObj[key]); */

const convertBytesToKB = (bytes) => Math.round(bytes / KILO_BYTES_PER_BYTE);

const UploadImage = () => {
  const axiosPrivate = useAxiosPrivate();
  const fileInputField = useRef(null);
  const [files, setFiles] = useState([]);
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(null);
  const [message, setMessage] = useState(null);

  const handleUploadBtnClick = () => {
    fileInputField.current.click();
  };

  const handleNewFileUpload = (e) => {
    const { files: newFiles } = e.target;
    setFiles([...files, ...newFiles]);
  };

  const removeFile = (fileName) => {
    const newArray = files.filter((file) => file.name !== fileName);
    setFiles([...newArray]);
  };

  const upload = async (e) => {
    e.preventDefault();
    const formData = new FormData();
    files.map(file =>
      formData.append('file', file)
    );

    try {
      setLoading(true);
      const res = await axiosPrivate({
        method: 'post',
        url: '/images',
        data: formData,
        headers: { 'Content-Type': 'multipart/form-data' }
      });
      setLoading(false);
      if (res.status === 200) {
        setFiles([]);
        setSuccess(true);
        setMessage('Images uploaded successfully');
      } else {
        setSuccess(false);
        setMessage('Error uploading images. Please try again.');
      }
      setTimeout(() => { setMessage(null); setSuccess(null); }, 2000);
    } catch (error) {
      setFiles([]);
      setSuccess(false);
      setMessage('Error uploading images. Please try again.');
      setLoading(false);
      setTimeout(() => { setMessage(null); setSuccess(null); }, 5000);
      console.log(error);
      return error;
    }
  };

  return (
    <>
      { loading
        ? <div className="loader-container">
          <div className="spinner"></div>
        </div>
        : <div className="background">
          <h3 className="fu-title">Upload here your new journey</h3>
          <FileUploadContainer className="fu-container">
            <DragDropText className="block">Drag and drop your files anywhere or</DragDropText>
            <UploadFileBtn className="block" type="button" onClick={handleUploadBtnClick}>
              <i className="fas fa-file-upload" />
              <span> Upload files</span>
            </UploadFileBtn>
            <FormField
              type="file"
              ref={fileInputField}
              onChange={handleNewFileUpload}
              title=""
              value=""
              multiple
            />
          </FileUploadContainer>
          {
            files.length === 0
              ? <br></br>
              : <FilePreviewContainer>
                <span>To Upload</span>
                <PreviewList>
                  {files.map((file) => {
                    return (
                      <PreviewContainer key={file.name}>
                        <div>
                          <ImagePreview
                            src={URL.createObjectURL(file)}
                          />
                          <FileMetaData isImageFile={true}>
                            <span>{file.name}</span>
                            <aside>
                              <span>{convertBytesToKB(file.size)} kb</span>
                              <RemoveFileIcon
                                className="fas fa-trash-alt"
                                onClick={() => removeFile(file.name)}
                              />
                            </aside>
                          </FileMetaData>
                        </div>
                      </PreviewContainer>
                    );
                  })}
                </PreviewList>
              </FilePreviewContainer>
          }
          {message &&
            <p className="message" style={{ backgroundColor: `${success ? 'green' : 'red'}` }}>{message}</p>}
          <button className="fu-button" onClick={upload}>Confirm Upload</button>

        </div>
      }
    </>
  );
};

export default UploadImage;
