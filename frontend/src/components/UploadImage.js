import React, { useRef, useState } from 'react';
import { uploadImages } from '../api/RemoteServices';
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
  const fileInputField = useRef(null);
  const [files, setFiles] = useState([]);

  const handleUploadBtnClick = () => {
    fileInputField.current.click();
  };

  /* const addNewFiles = (newFiles) => {
    for (let file of newFiles) {
      if (file.size < DEFAULT_MAX_FILE_SIZE_IN_BYTES) {
        files[file.name] = file;
      }
    }
    return { ...files };
  }; */

  const handleNewFileUpload = (e) => {
    const { files: newFiles } = e.target;
    setFiles([...files, ...newFiles]);
  };

  const removeFile = (fileName) => {
    const newArray = files.filter((file) => file.name !== fileName);
    setFiles([...newArray]);
    console.log(files.length);
  };

  const upload = async () => {
    const res = await uploadImages(files);
    // TODO deal with response
  };

  return (
    <div className="background">
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
        files.length == 0
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
      <button className="fu-button" onClick={upload}>Confirm Upload</button>

    </div>
  );
};

export default UploadImage;
