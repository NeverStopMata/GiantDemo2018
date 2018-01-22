package math

import (
//"fmt"
)

//求点到线段的距离，返回相交的位置（0~1），以及到相交点的距离
func PointToSegment(p *Vector3, segStart *Vector3, segEnd *Vector3) (cursor float32, dist float32) {
	diff := p.Substract(segStart)
	segDir := segEnd.Substract(segStart)
	cursor = diff.Dot(segDir)
	if cursor > 0 {
		segDisplacement := segDir.SqrMagnitude()
		if cursor >= segDisplacement {
			cursor = 1.0
			diff.DecreaseBy(segDir)
		} else {
			cursor /= segDisplacement
			diff.DecreaseBy(segDir.ScaleBy(cursor))
		}
	} else {
		//在起始点后面
		cursor = 0
	}
	dist = diff.Magnitude()
	return
}

//求点到直线的距离，返回到相交点的距离
func PointToLine(p *Vector3, lineP1, lineP2 *Vector3) (intesectPos *Vector3, dist float32) {
	vU := lineP2.Substract(lineP1)
	vU.NormalizeSelf()
	vAP := p.Substract(lineP1)
	vU.ScaleBy(vU.Dot(vAP))
	vV := vU.Add(lineP1)

	intesectPos = vV
	dist = GetDistance2D(vV, p)
	return
}

//将两个向量转化为四方向的枚举值
const (
	DIR_REL_FORWARD = iota
	DIR_REL_BACKWARD
	DIR_REL_LEFT
	DIR_REL_RIGHT
)

func PointRelationshipWithPosDir(basePos *Vector3, forwardDirection *Vector3, point *Vector3, offsetMargin float32) int {
	dist := GetDistance2DSqr(basePos, point)
	if LessEqualF32(dist, SqrF32(offsetMargin)) {
		return DIR_REL_FORWARD
	}
	dir := GetDirection2D(point, basePos)
	if dir.IsZero() {
		return DIR_REL_FORWARD
	}
	fd := forwardDirection.Normalize()
	dir.NormalizeSelf()
	dotResult := fd.Dot(dir)
	if dotResult > 0.707 {
		return DIR_REL_FORWARD
	} else if dotResult < -0.707 {
		return DIR_REL_BACKWARD
	}
	crossRet := fd.Cross(dir)
	if crossRet.Y > 0 {
		return DIR_REL_RIGHT
	}
	return DIR_REL_LEFT
}

//求解2元一次线性方程
// a1x + b1y = c1
// a2x + b2y = c2
func SolveLinearEquationWithTwoUnknowns(a1, b1, c1, a2, b2, c2 float32) (ok bool, x float32, y float32) {
	det := a1*b2 - b1*a2
	if IsZero(det) {
		x, y = 0, 0
		ok = false
		return
	}
	x = (b2*c1 - b1*c2) / det
	y = (a1*c2 - a2*c1) / det
	ok = true
	return
}

//求两根2D平面上的直线是否相交
func LineIntersect2D(s1, e1, s2, e2 *Vector3) (ok bool, intersectPoint *Vector3, cursor float32) {
	cursor = float32(0)
	ok, cursor, _ = SolveLinearEquationWithTwoUnknowns(e1.X-s1.X, s2.X-e2.X, s2.X-s1.X, e1.Z-s1.Z, s2.Z-e2.Z, s2.Z-s1.Z)
	intersectPoint = e1.Substract(s1).ScaleBy(cursor).IncreaseBy(s1)
	return
}

type EROConvertType int

const (
	EROConvertType_DIR_4 EROConvertType = iota + 1
	EROConvertType_DIR_8
	EROConvertType_DIR_16
)

func GetNormalizedOrientationV3(vec *Vector3, convertType EROConvertType) *Vector3 {
	if vec.IsZero() {
		return Vector3_0.Clone()
	}
	vRet := Vector3_0.Clone()
	fk := float32(0)
	fTanPID8 := float32(PI / 8)
	fTan3PID8 := fTanPID8 * 3
	fSqrt2 := SqrtF32(2)

	fTanPID16 := fTanPID8 / 2
	fTan3PID16 := fTanPID16 + (fTan3PID8-fTanPID8)/4 + fTanPID8/2
	fTan5PID16 := fTan3PID16 + (fTan3PID8-fTanPID8)/2
	fTan7PID16 := fTan5PID16 + (fTan3PID8-fTanPID8)/4 + (PI/2-fTan3PID8)/2
	if IsZero(vec.X) == false {
		fk = AbsF32(vec.Z) / AbsF32(vec.X)
	} else {
		if convertType == EROConvertType_DIR_4 {
			fk = 1
		} else if convertType == EROConvertType_DIR_16 {
			fk = TanF32(fTan7PID16)
		} else {
			fk = TanF32(fTan3PID8)
		}
	}
	if convertType == EROConvertType_DIR_4 {
		if fk < 1 && fk > -1 {
			if vec.X > 0 {
				vRet.X = 1
			} else {
				vRet.X = -1
			}
		} else {
			if vec.Z > 0 {
				vRet.Z = 1
			} else if vec.Z < 0 {
				vRet.Z = -1
			}
		}
	} else if convertType == EROConvertType_DIR_16 {
		fTanPID16 = TanF32(fTanPID16)
		fTan3PID16 = TanF32(fTan3PID16)
		fTan5PID16 = TanF32(fTan5PID16)
		fTan7PID16 = TanF32(fTan7PID16)
		if fk < fTanPID16 {
			if vec.X > 0 {
				vRet.X = 1
			} else {
				vRet.X = -1
			}
		} else if fk >= fTanPID16 && fk < fTan3PID16 {
			if vec.X > 0 && vec.Z > 0 {
				vRet.X = 1
				vRet.Z = fSqrt2 - 1
			} else if vec.X > 0 && vec.Z < 0 {
				vRet.X = 1
				vRet.Z = 1 - fSqrt2
			} else if vec.X < 0 && vec.Z > 0 {
				vRet.X = -1
				vRet.Z = fSqrt2 - 1
			} else if vec.X < 0 && vec.Z < 0 {
				vRet.X = -1
				vRet.Z = 1 - fSqrt2
			}
		} else if fk >= fTan3PID16 && fk < fTan5PID16 {
			if vec.X > 0 && vec.Z > 0 {
				vRet.X = 1
				vRet.Z = 1
			} else if vec.X > 0 && vec.Z < 0 {
				vRet.X = 1
				vRet.Z = -1
			} else if vec.X < 0 && vec.Z > 0 {
				vRet.X = -1
				vRet.Z = 1
			} else if vec.X < 0 && vec.Z < 0 {
				vRet.X = -1
				vRet.Z = -1
			}
		} else if fk >= fTan5PID16 && fk < fTan7PID16 {
			if vec.X > 0 && vec.Z > 0 {
				vRet.X = 1
				vRet.Z = fSqrt2 + 1
			} else if vec.X > 0 && vec.Z < 0 {
				vRet.X = 1
				vRet.Z = -1 - fSqrt2
			} else if vec.X < 0 && vec.Z > 0 {
				vRet.X = -1
				vRet.Z = fSqrt2 + 1
			} else if vec.X < 0 && vec.Z < 0 {
				vRet.X = -1
				vRet.Z = -1 - fSqrt2
			}
		} else if fk >= fTan7PID16 {
			if vec.Z > 0 {
				vRet.Z = 1
			} else if vec.Z < 0 {
				vRet.Z = -1
			}
		}
	} else if convertType == EROConvertType_DIR_8 {
		fTanPID8 = TanF32(fTanPID8)
		fTan3PID8 = TanF32(fTan3PID8)
		if fk < fTanPID8 {
			if vec.X > 0 {
				vRet.X = 1
			} else {
				vRet.X = -1
			}
		} else if fk >= fTanPID8 && fk < fTan3PID8 {
			if vec.X > 0 && vec.Z > 0 {
				vRet.X = 1
				vRet.Z = 1
			} else if vec.X > 0 && vec.Z < 0 {
				vRet.X = 1
				vRet.Z = -1
			} else if vec.X < 0 && vec.Z > 0 {
				vRet.X = -1
				vRet.Z = 1
			} else if vec.X < 0 && vec.Z < 0 {
				vRet.X = -1
				vRet.Z = -1
			}
		} else if fk >= fTan3PID8 {
			if vec.Z > 0 {
				vRet.Z = 1
			} else if vec.Z < 0 {
				vRet.Z = -1
			}
		}
	}
	vRet.NormalizeSelf()
	return vRet
}

type EDirection4 int

const (
	EDirection4_0 EDirection4 = iota
	EDirection4_1
	EDirection4_2
	EDirection4_3
	EDirection4_4
)

type EDirection8 int

const (
	EDirection8_0 EDirection8 = iota
	EDirection8_1
	EDirection8_2
	EDirection8_3
	EDirection8_4
	EDirection8_5
	EDirection8_6
	EDirection8_7
	EDirection8_8
)

type EDirection16 int

const (
	EDirection16_0 EDirection16 = iota
	EDirection16_1
	EDirection16_2
	EDirection16_3
	EDirection16_4
	EDirection16_5
	EDirection16_6
	EDirection16_7
	EDirection16_8
	EDirection16_9
	EDirection16_10
	EDirection16_11
	EDirection16_12
	EDirection16_13
	EDirection16_14
	EDirection16_15
	EDirection16_16
)

func GetNormalizedOrientation4Enum(outDir *Vector3, baseDir *Vector3) EDirection4 {
	return EDirection4(getNormalizedOrientationEnum(outDir, baseDir, EROConvertType_DIR_4))
}
func GetNormalizedOrientation8Enum(outDir *Vector3, baseDir *Vector3) EDirection8 {
	return EDirection8(getNormalizedOrientationEnum(outDir, baseDir, EROConvertType_DIR_8))
}
func GetNormalizedOrientation16Enum(outDir *Vector3, baseDir *Vector3) EDirection16 {
	return EDirection16(getNormalizedOrientationEnum(outDir, baseDir, EROConvertType_DIR_16))
}
func getNormalizedOrientationEnum(outDir *Vector3, baseDir *Vector3, convertType EROConvertType) int {
	if outDir.IsZero() {
		return 0
	}
	if baseDir.IsZero() {
		return 0
	}
	dirCount := 0
	if convertType == EROConvertType_DIR_4 {
		dirCount = 4
	} else if convertType == EROConvertType_DIR_8 {
		dirCount = 8
	} else if convertType == EROConvertType_DIR_16 {
		dirCount = 16
	} else {
		return 0
	}
	normalizedOutDir := outDir.Normalize()
	normalizedBaseDir := baseDir.Normalize()
	angle := AngleBetween2DWithSign(normalizedOutDir, normalizedBaseDir)
	angle = NormalizeAngleZeroToTwoPI(angle)
	angleF32 := angle.ToFloat32()
	//fmt.Println(angleF32, dirCount)
	angleStep := PI * 2 / float32(dirCount)
	for i := 0; i <= dirCount; i++ {
		dAngle := angleF32 - angleStep*float32(i)
		//fmt.Println(dAngle)
		if AbsF32(dAngle) <= angleStep*0.5 {
			if i == dirCount {
				return 1
			}
			return i + 1
		}
	}
	return 0
}
func GetDirectionByDirEnum4(dirEnum EDirection4) *Vector3 {
	v := new(Vector3).Create()
	if dirEnum == EDirection4_0 {
		return v
	}
	step := float32(PI / 2)
	deltaEnum := int(dirEnum) - int(EDirection4_1)
	v.X = CosF32(float32(deltaEnum) * step)
	v.Z = -SinF32(float32(deltaEnum) * step)
	v.NormalizeSelf()
	return v
}
func GetDirectionByDirEnum8(dirEnum EDirection8) *Vector3 {
	v := new(Vector3).Create()
	if dirEnum == EDirection8_0 {
		return v
	}
	step := float32(PI / 4)
	deltaEnum := int(dirEnum) - int(EDirection8_1)
	v.X = CosF32(float32(deltaEnum) * step)
	v.Z = -SinF32(float32(deltaEnum) * step)
	v.NormalizeSelf()
	return v
}
func GetDirectionByDirEnum16(dirEnum EDirection16) *Vector3 {
	v := new(Vector3).Create()
	if dirEnum == EDirection16_0 {
		return v
	}
	step := float32(PI / 8)
	deltaEnum := int(dirEnum) - int(EDirection16_1)
	v.X = CosF32(float32(deltaEnum) * step)
	v.Z = -SinF32(float32(deltaEnum) * step)
	v.NormalizeSelf()
	return v
}
