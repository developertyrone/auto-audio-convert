#!/bin/bash
# Integration test for auto-audio-convert
# Note: Requires ffmpeg installed

set -e

TEST_DIR="/tmp/auto-audio-test-$$"
BINARY="./auto-audio-convert"

echo "🧪 Running integration test..."
echo ""

# Setup test directory
mkdir -p "$TEST_DIR/music/subfolder1/subfolder2"

echo "📁 Creating test structure in $TEST_DIR..."

# Create fake audio files (in real scenario, ffmpeg would process these)
cat << 'EOF' > "$TEST_DIR/music/song1.flac"
fake flac file 1
EOF

cat << 'EOF' > "$TEST_DIR/music/song2.flac"
fake flac file 2
EOF

cat << 'EOF' > "$TEST_DIR/music/subfolder1/song3.flac"
fake flac file 3
EOF

cat << 'EOF' > "$TEST_DIR/music/subfolder1/subfolder2/song4.flac"
fake flac file 4
EOF

# Create a pre-existing MP3 to test skip logic
cat << 'EOF' > "$TEST_DIR/music/song1.mp3"
already converted
EOF

echo "✅ Test files created:"
tree "$TEST_DIR" 2>/dev/null || find "$TEST_DIR" -type f

echo ""
echo "🔍 Testing scanner (without ffmpeg)..."
$BINARY --source="$TEST_DIR/music" --from=flac --to=mp3 2>&1 | head -10

echo ""
echo "📊 Expected behavior:"
echo "  - Should find 4 .flac files"
echo "  - Should skip song1.flac (song1.mp3 exists)"
echo "  - Should fail conversion (no ffmpeg) but show correct detection"

echo ""
echo "🧹 Cleaning up..."
rm -rf "$TEST_DIR"

echo "✅ Test complete!"
