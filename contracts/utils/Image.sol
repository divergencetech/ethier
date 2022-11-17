// SPDX-License-Identifier: MIT
// Copyright (c) 2022 the ethier authors (github.com/divergencetech/ethier)
pragma solidity >=0.8.16 <0.9.0;

import {DynamicBuffer} from "./DynamicBuffer.sol";

/**
 * @notice The rectangle defining a pixel frame in relation to a global
 * index coordinate system.
 */
struct Rectangle {
    uint8 xMin;
    uint8 yMin;
    uint8 xMax;
    uint8 yMax;
}

/**
 * @notice Utilities library to work with raw pixel data.
 * @dev The code assumes (32)24-bit (A)BGR pixel encoding.
 * @dev Frames without any explicit rectangle information are assumed to start
 * at the coordinate origin `xMin = yMin = 0`.
 */
//solhint-disable no-empty-blocks
library Image {
    using DynamicBuffer for bytes;

    /**
     * @notice Fills a pixel buffer with a given RGB color.
     */
    function fill(bytes memory bgrPixels, uint24 rgb) internal pure {
        assembly {
            let bgr := shl(
                // Pushing the BGR tripplet all the way to the left 256 - 24
                232,
                or(
                    and(0x00FF00, rgb),
                    or(shl(16, and(0xFF, rgb)), and(0xFF, shr(16, rgb)))
                )
            )

            bgr := or(bgr, shr(24, bgr))
            {
                let bgr2 := bgr
                bgr := or(bgr, shr(48, bgr))
                bgr := or(bgr, shr(96, bgr))
                bgr := or(bgr, shr(192, bgr2))
            }

            let bufPtr := add(bgrPixels, 0x20)
            let bufPtrEnd := add(bufPtr, mload(bgrPixels))
            for {

            } 1 {

            } {
                // Stopping if we reached the end of the block.
                if iszero(lt(add(bufPtr, 32), bufPtrEnd)) {
                    break
                }

                mstore(bufPtr, bgr)
                bufPtr := add(bufPtr, 30)
            }

            let mask := shr(
                shl(3, sub(bufPtrEnd, bufPtr)),
                0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
            )

            mstore(bufPtr, or(and(bgr, not(mask)), and(mload(bufPtr), mask)))
        }
    }

    /**
     * @notice Mask the location of the vectorised BGR channels.
     */
    uint256 private constant _VECTORISATION_MASK =
        0xFF0000000000000000FF0000000000000000FF;

    /**
     * @notice Blends two pixels depending on the alpha channel of the latter.
     * @dev An accuracy-focused algorithm that removes bias across color
     * channels. See also https://stackoverflow.com/a/1230272
     * @param bgr BGR encoded pixel.
     * @param abgr ABGR encoded pixel with alpha channel.
     */
    function alphaBlend(uint24 bgr, uint32 abgr)
        internal
        pure
        returns (uint24 res)
    {
        assembly {
            let a := shr(24, abgr)
            let na := sub(0xff, a)

            // Spacing the color channel values across the 256 bit word.
            // | 0 (13B) | R (1B) | 0 (8B) | G (1B) | 0 (8B) | R (1B) |
            // This allows all channels to be blended in a single operation.
            bgr := and(
                or(or(shl(128, bgr), shl(64, bgr)), bgr),
                _VECTORISATION_MASK
            )

            abgr := and(
                or(or(shl(128, abgr), shl(64, abgr)), abgr),
                _VECTORISATION_MASK
            )

            // h = alpha * fg + (255 - alpha) * bg + 128
            let h := add(
                add(mul(a, abgr), mul(na, bgr)),
                // Adds 0x80 to each value
                0x80000000000000000080000000000000000080
            )

            // h = ((h >> 8) + h) >> 8
            h := and(
                shr(8, add(shr(8, h), h)),
                // Bit cleaning
                _VECTORISATION_MASK
            )

            res := or(or(shr(128, h), shr(64, h)), h)
        }
    }

    /**
     * @notice Blends a background frame with foreground one depending on the
     * alpha channel of the latter.
     * @param bgBgr BGR encoded pixel frame (background)
     * @param fgAbgr ABGR encoded pixel frame with alpha channel
     * (foreground)
     * @param width of the background frame
     * @param rect The frame rectangle (coordinates) of the foreground
     */
    function alphaBlend(
        bytes memory bgBgr,
        bytes memory fgAbgr,
        uint256 width,
        Rectangle memory rect
    ) internal pure {
        uint256 fgStride = (rect.xMax - rect.xMin) * 4;
        uint256 bgStride = width * 3;

        uint256 fgCursor;
        uint256 bgCursor;
        assembly {
            fgCursor := add(fgAbgr, 0x20)
            bgCursor := add(bgBgr, 0x20)
        }

        // Adding the offset to the lower left corner of the foreground frame
        bgCursor += rect.xMin * 3 + rect.yMin * bgStride;

        // The background pointer jump going from the end of one row in the
        // foreground frame to the start of the next one.
        uint256 rowJump = bgStride - (rect.xMax - rect.xMin) * 3;

        assembly {
            // This computation kernel has been taken and inlined from
            // `alphaBlend(uint24 bgr, uint32 abgr)` for efficiency.
            function alphaBlend(bgrPtr, abgrPtr) {
                let buf := mload(bgrPtr)
                let bgr := shr(232, buf)
                let abgr := shr(224, mload(abgrPtr))

                let a := shr(24, abgr)
                let na := sub(0xff, a)

                // Spacing the color channel values across the 256 bit word.
                // | 0 (13B) | R (1B) | 0 (8B) | G (1B) | 0 (8B) | R (1B) |
                // This allows all channels to be blended in a single operation.
                bgr := and(
                    or(or(shl(128, bgr), shl(64, bgr)), bgr),
                    _VECTORISATION_MASK
                )

                abgr := and(
                    or(or(shl(128, abgr), shl(64, abgr)), abgr),
                    _VECTORISATION_MASK
                )

                // h = alpha * fg + (255 - alpha) * bg + 128
                let h := add(
                    add(mul(a, abgr), mul(na, bgr)),
                    // Adds 0x80 to each value
                    0x80000000000000000080000000000000000080
                )

                // h = ((h >> 8) + h) >> 8
                h := and(
                    shr(8, add(shr(8, h), h)),
                    // Bit cleaning
                    _VECTORISATION_MASK
                )

                let res := or(or(shr(128, h), shr(64, h)), h)

                mstore(
                    bgrPtr,
                    or(
                        shl(232, res),
                        and(
                            buf,
                            0x000000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
                        )
                    )
                )
            }

            // Looping over the foreground frame
            let fgEnd := add(fgCursor, mload(fgAbgr))
            let fgIdx := 0

            for {

            } 1 {

            } {
                // Stopping if we reached the end of the foreground frame.
                if iszero(lt(fgCursor, fgEnd)) {
                    break
                }

                alphaBlend(bgCursor, fgCursor)

                fgIdx := add(fgIdx, 4)
                fgCursor := add(fgCursor, 4)
                bgCursor := add(bgCursor, 3)

                // If we are switching rows in the foreground frame, we have to
                // make a larger jump for the background cursor.
                if iszero(mod(fgIdx, fgStride)) {
                    bgCursor := add(bgCursor, rowJump)
                }
            }
        }
    }

    /**
     * @notice Scales a pixel frame.
     * @param bgr BGR encoded pixel frame.
     * @param width of the frame.
     * @param pixelSize The number of bytes in a pixel. e.g. 4 for ABGR.
     * @param scalingFactor The scaling factor.
     */
    function scale(
        bytes memory bgr,
        uint256 width,
        uint256 pixelSize,
        uint256 scalingFactor
    ) internal pure returns (bytes memory) {
        bytes memory buffer = DynamicBuffer.allocate(
            bgr.length * scalingFactor * scalingFactor
        );
        appendSafeScaled(buffer, bgr, width, pixelSize, scalingFactor);
        return buffer;
    }

    /**
     * @notice Scales a pixel frame and appends the rescaled data to a given
     * buffer.
     * @dev This routine is compatible with ethier's `DynamicBuffer`.
     * @param bgr BGR encoded pixel frame.
     * @param width of the frame.
     * @param pixelSize The number of bytes in a pixel. e.g. 4 for ABGR.
     * @param scalingFactor The scaling factor.
     */
    function appendSafeScaled(
        bytes memory buffer,
        bytes memory bgr,
        uint256 width,
        uint256 pixelSize,
        uint256 scalingFactor
    ) internal pure {
        buffer.checkOverflow(bgr.length * scalingFactor * scalingFactor);

        assembly {
            /**
             * @notice Fills a 2D block in memory by repeating linear chunks
             * e.g.
             * | ..................................|
             * | .... | chunk | chunk | tail | ... |
             * | .... | chunk | chunk | tail | ... |
             * | .... | chunk | chunk | tail | ... |
             * | ..................................|
             * where tail is a broken chunk
             * @param bufPtr The memory pointer to the upper left corner of the
             * block
             * @param bufStride The buffer stride, i.e. the number of bytes that
             * need to be added to get from one row of the buffer to the next
             * without changing the column (aka. the buffer width)
             * @param blockWidth The number of columns covered by the block
             * @param blockHeight The number of rows covered by the block
             * @param chunk The bytes that will be used to fill the block
             * (single word, i.e. max 32 bytes). Big endian, i.e.
             * chunk[:chunkSize] will be used.
             * @param chunkSize The size of the chunk. See above.
             * @param tailMask Mask the bits of the chunk that have to be used
             * for the tail of the block.
             */
            function writeBlock(
                bufPtr,
                bufStride,
                blockWidth,
                blockHeight,
                chunk,
                chunkSize,
                tailMask
            ) {
                // The pointer to the lower right corner of the block
                let bufPtrEnd := add(bufPtr, mul(bufStride, blockHeight))

                // Row loop
                for {

                } 1 {

                } {
                    // Stopping if we reached the end of the block.
                    if iszero(lt(bufPtr, bufPtrEnd)) {
                        break
                    }

                    let rowPtr := bufPtr

                    // Column loop
                    // We are going to write chunks as full words for efficiency.
                    // This might result in out-of-bound writes at the row tail
                    // which will thus need special treatment (masking).
                    for {
                        // Stopping a word before the end of the chunk row to
                        // treat the tail separately.
                        let rowEnd := sub(add(rowPtr, blockWidth), 0x20)
                    } 1 {

                    } {
                        if lt(rowEnd, rowPtr) {
                            break
                        }
                        mstore(rowPtr, chunk)
                        rowPtr := add(rowPtr, chunkSize)
                    }

                    // Since writing a full word would affect memory outside of
                    // the block we load the current content and mix it with the
                    // tail data.
                    mstore(
                        rowPtr,
                        or(
                            and(chunk, not(tailMask)),
                            and(mload(rowPtr), tailMask)
                        )
                    )
                    bufPtr := add(bufPtr, bufStride)
                }
            }

            let dataPtr := add(bgr, 0x20)
            let dataPtrEnd := add(dataPtr, mload(bgr))
            let dataIdx := sub(dataPtr, add(bgr, 0x20))
            let dataStride := mul(width, pixelSize)

            let chunkSize := mul(div(32, pixelSize), pixelSize)
            let blockWidth := mul(pixelSize, scalingFactor)

            let bufPtr := add(add(buffer, 0x20), mload(buffer))
            let bufStride := mul(dataStride, scalingFactor)

            let pixelMask := not(
                shr(
                    shl(3, pixelSize), // * 8
                    0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
                )
            )

            // Binary mask for the tail of the block (i.e. the last
            // chunk that will only be partially written)
            let tailMask := shr(
                shl(3, mod(blockWidth, chunkSize)), // * 8
                0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
            )

            // Loop over all data pixels
            for {

            } 1 {

            } {
                if iszero(lt(dataPtr, dataPtrEnd)) {
                    break
                }

                // Building the "chunk" by repeatedly appending the pixel data
                // until a 256bit stack word is full
                let chunk := 0
                for {
                    let pixel := and(mload(dataPtr), pixelMask)
                    let size := 0
                    let shift := mul(pixelSize, 8)
                } 1 {

                } {
                    if iszero(lt(size, chunkSize)) {
                        break
                    }

                    chunk := or(chunk, pixel)
                    pixel := shr(shift, pixel)
                    size := add(size, pixelSize)
                }

                // Fill the block with pixel data
                writeBlock(
                    bufPtr,
                    bufStride,
                    blockWidth,
                    scalingFactor,
                    chunk,
                    chunkSize,
                    tailMask
                )

                dataIdx := add(dataIdx, pixelSize)
                dataPtr := add(dataPtr, pixelSize)
                bufPtr := add(bufPtr, blockWidth)

                // If we are switching rows in the block, we have to make a
                // larger jump for the buffer cursor.
                if iszero(mod(dataIdx, dataStride)) {
                    bufPtr := add(
                        bufPtr,
                        mul(sub(bufStride, dataStride), scalingFactor)
                    )
                }
            }

            // Update the length of the buffer
            mstore(
                buffer,
                add(
                    mload(buffer),
                    mul(mload(bgr), mul(scalingFactor, scalingFactor))
                )
            )
        }
    }
}
//solhint-enable no-empty-blocks
