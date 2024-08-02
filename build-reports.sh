#!/bin/sh
cat ${REPORTS_DIR}/reports.mix | awk 'BEGIN{pr=0} { if($0 ~ /Begin cart report/)   {pr=1; next}; if($0 ~ /End cart report/)   {exit}; if(pr==1) print $0}' > ${REPORTS_DIR}/test-report-cart.xml
cat ${REPORTS_DIR}/reports.mix | awk 'BEGIN{pr=0} { if($0 ~ /Begin loms report/)   {pr=1; next}; if($0 ~ /End loms report/)   {exit}; if(pr==1) print $0}' > ${REPORTS_DIR}/test-report-loms.xml
cat ${REPORTS_DIR}/reports.mix | awk 'BEGIN{pr=0} { if($0 ~ /Begin cart coverage/) {pr=1; next}; if($0 ~ /End cart coverage/) {exit}; if(pr==1) print $0}' > ${REPORTS_DIR}/coverage-cart.xml
cat ${REPORTS_DIR}/reports.mix | awk 'BEGIN{pr=0} { if($0 ~ /Begin loms coverage/) {pr=1; next}; if($0 ~ /End loms coverage/) {exit}; if(pr==1) print $0}' > ${REPORTS_DIR}/coverage-loms.xml
