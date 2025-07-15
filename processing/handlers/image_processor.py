import json
from typing import Optional
from PIL import Image
from io import BytesIO
from pathlib import Path
from pydantic import BaseModel


class ResizeOperation(BaseModel):
    width: int
    height: int


class ImageOperation(BaseModel):
    resize: Optional[ResizeOperation]
    grayscale: Optional[bool]


class ImageProcessingPayload(BaseModel):
    input_path: str
    output_path: str
    operations: list[ImageOperation]


class ImageProcessor:
    @staticmethod
    def process_image(payload: ImageProcessingPayload):
        """Process image based on task payload"""
        try:
            with Image.open(payload['input_path']) as img:
                for operation in payload['operations']:
                    img = ImageProcessor._apply_operation(img, operation)

                output_format = Path(
                    payload['output_path']).suffix[1:].upper()
                img.save(payload['output_path'], format=output_format)

            return True, "Image processed successfully"
        except Exception as e:
            return False, str(e)

    @staticmethod
    def _apply_operation(img, operation):
        """Apply single image operation"""
        if 'resize' in operation:
            return img.resize(
                (operation['resize']['width'],
                 operation['resize']['height']))
        elif 'grayscale' in operation:
            return img.convert('L')
        elif 'rotate' in operation:
            return img.rotate(operation['rotate'])
        elif 'crop' in operation:
            box = (
                operation['crop']['left'],
                operation['crop']['top'],
                operation['crop']['right'],
                operation['crop']['bottom']
            )
            return img.crop(box)
        return img


def handle_task(payload: ImageProcessingPayload):
    payload = json.loads(payload)
    success, message = ImageProcessor.process_image(payload)

    return {
        'success': success,
        'message': message,
        'output_path': payload['output_path']
    }
