import React, { useRef, useState } from "react";
import { uploadImages } from "../services/RemoteServices"
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
  RemoveFileIcon,
  InputLabel
} from "../styles/file-upload";

const KILO_BYTES_PER_BYTE = 1000;
const DEFAULT_MAX_FILE_SIZE_IN_BYTES = 500000;

const convertNestedObjectToArray = (nestedObj) =>
  Object.keys(nestedObj).map((key) => nestedObj[key]);

const convertBytesToKB = (bytes) => Math.round(bytes / KILO_BYTES_PER_BYTE);

const UploadImage = () => {
  const fileInputField = useRef(null);
  const [files, setFiles] = useState([]);

  const handleUploadBtnClick = () => {
    fileInputField.current.click();
  };

  /*const addNewFiles = (newFiles) => {
    for (let file of newFiles) {
      if (file.size < DEFAULT_MAX_FILE_SIZE_IN_BYTES) {
        files[file.name] = file;
      }
    }
    return { ...files };
  };*/

  const handleNewFileUpload = (e) => {
    const { files: newFiles } = e.target;
    console.log(newFiles)
    let newFile = e.target.files[0]
    /*files.map((file) => {
      if(file.name == newFile.name) {
        //TODO raise error. Duplicate file with same name
      } 
    })*/

    setFiles([...files, ...newFiles]);
    console.log(files.length)
  };

  const removeFile = (fileName) => {
    let newArray = files.filter((file) => file.name !== fileName)
    setFiles([ ...newArray ]);
    console.log(files.length)
  };

  const upload = async () => {
    const res = await uploadImages(files)
    //TODO deal with response
  };



  return (
    <>
      <FileUploadContainer>
        <InputLabel>Upload here your new journey</InputLabel>
        <DragDropText>Drag and drop your files anywhere or</DragDropText>
        <UploadFileBtn type="button" onClick={handleUploadBtnClick}>
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
      <FilePreviewContainer>
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
      <button onClick={upload}>Confirm Upload</button>
    </>
  );
};

export default UploadImage;