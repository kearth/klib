#!/usr/bin/python3
# -*- coding: UTF-8 -*-

# 图片二值化
from PIL import Image
import sys
import os.path

fileEnd = ('.bmp', '.png', '.jpg', '.jpeg', '.tif', '.tiff')
argc = len(sys.argv)
if  argc == 2 or argc == 3 :
    imgFile = sys.argv[1]
    if imgFile.lower().endswith(fileEnd):
        img = Image.open(imgFile)
        # 模式L为灰色图像，它的每个像素用8个bit表示，0表示黑，255表示白，其他数字表示不同的灰度。
        Img = img.convert('L')
        if argc == 3:
            imgName = sys.argv[2]
        else:
            imgName = imgFile.split('.')[0] + "_gray"
        Img.save(imgName + ".jpg")
    else:
        print("图片格式不正确")
else :
    if os.path.isfile("test.jpg"):
        img = Image.open("test.jpg")
        Img = img.convert('L')
        Img.save("test_gray.jpg")
    else :
        print("参数必须是2个，第一个是本文件，第二个是图片")
